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

package capmc

// We goofed in the original CAPMC API when we didn't put a version
// number in the URL. These constants define paths with and without
// the version number for backward compatibility with published
// documentation. Newer code should use the v1 paths.

// The Shasta implementation of the Cascade CAPMC APIs
const (
	ComputeNodeControlV1   = "/capmc/v1/cnctl"
	HealthV1               = "/capmc/v1/health"
	LivenessV1             = "/capmc/v1/liveness"
	PowerCapCapabilitiesV1 = "/capmc/v1/get_power_cap_capabilities"
	PowerCapGetV1          = "/capmc/v1/get_power_cap"
	PowerCapSetV1          = "/capmc/v1/set_power_cap"
	ReadinessV1            = "/capmc/v1/readiness"
	XnameOffV1             = "/capmc/v1/xname_off"
	XnameOnV1              = "/capmc/v1/xname_on"
	XnameReinitV1          = "/capmc/v1/xname_reinit"
	XnameStatusV1          = "/capmc/v1/get_xname_status"
)
