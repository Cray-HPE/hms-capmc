// MIT License
//
// (C) Copyright [2019, 2021] Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//

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
