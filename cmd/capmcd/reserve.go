// Copyright  2020 Hewlett Packard Enterprise Development LP
//
// This file contains the CAPMC to hardware state manager (HSM) interfaces
//

package main

func (d *CapmcD) reserveComponents(xnames []string, cmd string) error {
	if !d.reservationsEnabled {
		return nil
	}

	targetMap := make(map[string]bool)
	var targetedXnames []string
	var descendants []string
	var err error

	//start with just the original xnames
	for _, xname := range xnames {
		targetMap[xname] = true
	}

	//if the cmd matches then get its descendants;
	switch cmd {
	case bmcCmdPowerOff, bmcCmdPowerForceOff,
		bmcCmdPowerRestart, bmcCmdPowerForceRestart:
		// Powering off a component will impact all of the components below it,
		// they also need to have their reservations set
		squery := HSMQuery{
			ComponentIDs: xnames,
			States:       []string{"!Empty"},
		}
		descendants, err = d.GenerateXnameDescendantList(squery)
		if err != nil {
			return err
		}
	}

	//add the descendants
	for _, child := range descendants {
		targetMap[child] = true
	}

	//create a final list of targets
	for xname, _ := range targetMap {
		targetedXnames = append(targetedXnames, xname)
	}

	err = d.reservation.Aquire(targetedXnames)
	return err
}

func (d *CapmcD) releaseComponents(xnames []string) error {
	if !d.reservationsEnabled {
		return nil
	}
	var clearList []string
	var error error
	for _, xname := range xnames {
		if d.reservation.Check([]string{xname}) == true {
			clearList = append(clearList, xname)
		}
	}
	error = d.reservation.Release(clearList)
	return error
}

// If there are any active reservations, remove them all
func (d *CapmcD) removeAllActiveReservations() {
	if !d.reservationsEnabled {
		return
	}

	currentReservations := d.reservation.Status()
	var targetedXnames []string
	for xname, _ := range currentReservations {
		targetedXnames = append(targetedXnames, xname)
	}
	_ = d.reservation.Release(targetedXnames)
}
