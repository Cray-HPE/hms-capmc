// Copyright 2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.

package tsdb

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

// AdjustTimestampsForHystersisAndMinimumWindow - will adjust the start time and end time to account for the hysteresis and the window
// If both times are before NOW - Hystersis; then do nothing, assuming minimum period is met
// If Both times are GREATER than NOW - Hystersis; then slide back in time so that both values are less than Now-Hyst; keep the duration equal, assuming the window is big enough
// If a time falls in between Hyst and Now; then slide the whole time window back in time, so that the End is < Now-Hyst; keeping the duration equal, assuming the window is big enough.
func AdjustTimestampsForHystersisAndMinimumWindow(startTime time.Time, endTime time.Time, rightNow time.Time, hysteresis time.Duration, window time.Duration) (actualStartTime time.Time, actualEndTime time.Time, err error) {
	log.Tracef("RECEIVED: start: %s end: %s now: %s hysteresis(s): %d window(s): %d ",
		startTime, endTime, rightNow, hysteresis/time.Second, window/time.Second)

	// If the right tail is greater than the Hystersis offset; SLIDE to the left by the overshot
	// This math LOOKS weird; but its actually:
	// if (endTime  - (NOW() - 15s))  > 0
	hysteresisTime := rightNow.Add(hysteresis)
	log.Tracef("hysteresisTime: %s ",
		hysteresisTime)
	overShot := endTime.Sub(hysteresisTime)
	log.Tracef("OVRSHOT(s): %d ",
		overShot/time.Second)

	if overShot > 0 {
		endTime = endTime.Add(overShot * -1)
		startTime = startTime.Add(overShot * -1)
	}

	if endTime.Sub(startTime) < window {
		startTime = endTime.Add(window * -1)
	}

	actualStartTime = startTime
	actualEndTime = endTime

	log.Tracef("RETURN: start: %s end: %s ",
		startTime, endTime)
	return

}

// ValidateTimeBoundNodeRequest -  validates that a TimeBoundNodeRequest follows the rules
// Immediate rules: Both startTime/endTime MUST be set && nodes list is NOT empty
// LONGTERM rules:
// Rules - Temporal : either jobID is set or BOTH startTime and endTime are set (and not 0 case)
// Rules - Geographical : either jobID is set or the Nodes list is not empty
// ONE RULE from EACH must be true!
func ValidateTimeBoundNodeRequest(tbnr TimeBoundNodeRequest) (valid bool, err error) {
	if (!tbnr.StartTime.IsZero() && !tbnr.EndTime.IsZero()) &&
		( len(tbnr.Nodes) > 0) {
		return true, nil
	}
	return false, errors.New(InvalidArguments)
}
