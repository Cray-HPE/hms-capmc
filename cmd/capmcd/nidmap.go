// Copyright 2017-2019 Cray Inc. All Rights Reserved.
//

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"stash.us.cray.com/HMS/hms-capmc/internal/capmc"
)

// newNidMapError creates a new NodeInfo structure initialized as an error
// response for the get nid map API.
func newNidInfoError(nid, ecode int, emsg string) *capmc.NidInfo {
	return &capmc.NidInfo{Nid: nid, E: ecode, ErrMsg: emsg}
}

// doNidMap handles the get_nid_map API
func (d *CapmcD) doNidMap(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		sendJsonError(w, http.StatusMethodNotAllowed,
			fmt.Sprintf("(%s) Not Allowed", r.Method))
		return
	}

	var args capmc.NidlistRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&args)
	if err != nil {
		if err == io.EOF {
			sendJsonError(w, http.StatusBadRequest, "no request")
		} else {
			sendJsonError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	var query = HSMQuery{
		States:  []string{"!Empty"},
		Enabled: []bool{true},
	}

	// The incoming NID list could be invalid. Do simple validation
	// before contacting Hardware State Manager.
	if len(args.Nids) > 0 {
		var invalidNIDs []int
		// Duplicate NIDs aren't an error in the Cascade CAPMC API.
		query.NIDs, invalidNIDs = validateNIDs(true, args.Nids)
		if len(invalidNIDs) > 0 {
			sendJsonError(w, http.StatusBadRequest,
				fmt.Sprintf("invalid nids: %v", invalidNIDs))
			return
		}
	}

	var data capmc.GetNidMapResponse

	nidInfo, err := d.GetNidInfo(query)
	if err != nil {
		var nidError *InvalidNIDsError

		if errors.As(err, &nidError) {
			for _, nid := range nidError.NIDs {
				data.Nids = append(data.Nids,
					newNidInfoError(nid, 22, "Undefined NID"))
			}
			data.E = 22 // EINVAL
			data.ErrMsg = "Invalid argument"
		} else {
			log.Printf("Error: %s", err)
			sendJsonError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// This shouldn't happen on a properly configured system
		// but it isn't necessarily an error.
		if nidInfo == nil {
			nidInfo = make([]*capmc.NidInfo, 0)
		}

		data.Nids = nidInfo
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
