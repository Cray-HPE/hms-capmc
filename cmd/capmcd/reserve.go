/*
 * MIT License
 *
 * (C) Copyright [2019-2021] Hewlett Packard Enterprise Development LP
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

package main

func (d *CapmcD) reserveComponents(xnames []string, cmd string) ([]string, error) {
	if !d.reservationsEnabled {
		return nil, nil
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
			return nil, err
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
	return targetedXnames, err
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
