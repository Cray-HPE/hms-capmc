# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!--
Guiding Principles:
* Changelogs are for humans, not machines.
* There should be an entry for every single version.
* The same types of changes should be grouped.
* Versions and sections should be linkable.
* The latest version comes first.
* The release date of each version is displayed.
* Mention whether you follow Semantic Versioning.

Types of changes:
Added - for new features
Changed - for changes in existing functionality
Deprecated - for soon-to-be removed features
Removed - for now removed features
Fixed - for any bug fixes
Security - in case of vulnerabilities
-->

## [1.25.6] - 2021-08-09

### Changed

- Added GitHub configuration files.

## [1.25.5] - 2021-07-27

### Changed

- Changed Stash to GitHub

## [1.25.4] - 2021-07-20

### Added

- Update image name and docker args for build

## [1.25.3] - 2021-07-19

### Added

- Conversion to github

## [1.25.2] - 2021-07-12

### Security

- CASMHMS-4933 - Updated base container images for security updates.

## [1.25.1] - 2021-06-21

### Changes

- Fixed filter response for On/Off when using Redfish

## [1.25.0] - 2021-06-07

### Changed

- Created release branch for CSM 1.2

## [1.24.0] - 2021-06-07

### Changed

- Created release branch for CSM 1.1

## [1.23.6] - 2021-05-04

### Changed

- Updated docker-compose files to pull images from Artifactory instead of DTR.

## [1.23.5] - 2021-04-22

### Fixed

- Chunked calls to HSM due to URL limits

## [1.23.4] - 2021-04-16

### Changed

- Updated Dockerfiles to pull base images from Artifactory instead of DTR.

## [1.23.3] - 2021-04-07

### Fixed

- Certain cases allowed parentless reservations preventing power operations

## [1.23.2] - 2021-03-15

### Changed

- Set CPU resource limits higher than the default

## [1.23.1] - 2021-02-05

### Changed

- Added User-Agent headers to outbound HTTP requests.

## [1.23.0] - 2021-02-05

### Changed

- Updated vendor packages and licensing

## [1.22.0] - 2021-02-03

### Changed
- Updated vendor packages

## [1.21.0] - 2021-01-14

### Changed

- Updated license file.


## [1.20.0] - 2020-12-11

### Changed
- CASMHMS-4279 - Fix setting power cap on GB nodes with newer firmware

## [1.19.11] - 2020-12-09

### Changed
- CASMHMS-4269 - Updated to support GB physical context change

## [1.19.10] - 2020-11-17

### Added
- CASMHMS-3781 - Added support for power control on Rosetta switches

## [1.19.9] - 2020-11-16

- CASMHMS-4214 - Added final CA bundle configmap handling to Helm chart.

## [1.19.8] - 2020-11-06

### Added
- CASMHMS-4124 - Added support for power capping HPE Proliant DL servers

## [1.19.7] - 2020-11-04

### Added
- CASMHMS-4123 - Support for power readings from HPE Proliant servers

## [1.19.6] - 2020-10-30

### Security
- CASMHMS-4105 - Updated base Golang Alpine image to resolve libcrypto vulnerability.

## [1.19.5] - 2020-10-29

Converted to v2 HSM.

## [1.19.4] - 2020-10-26

Added TLS verification capability to Redfish HTTP transports.

## [1.19.3] - 2020-10-16

### Changed
- CASMHMS-3939 - Don't allow power operations on disabled components

## [1.19.2] - 2020-10-01

### Security
- CASMHMS-4065 - Update base image to alpine-3.12.

## [1.19.1] - 2020-09-23

### Added
- CASMHMS-3962 - Support for HPE PushPowerButton

## [1.19.0] - 2020-09-15

### Changed
- Updates for new base charts, Helm v1/Loftsman v1, and support for Helm v3

## [1.18.5] - 2020-09-10

### Changed
- CASMHMS-3950 - Fixed issue where CAPMC was unable to get power status from Redfish Outlets.  
- Fixed an issue where CAPMC was unable to power off PDU Outlets.

## [1.18.4] - 2020-08-18

### Added 
- CASMHMS-3809 - Added support for the CabinetPDUPowerConnector HMSType

## [1.18.3] - 2020-07-15

### Changed 
- CASMHMS-3771 - Increased the HTTP client timeout from 20 seconds to 3 minutes.

## [1.18.2] - 2020-07-09

### Changed
- CASMHMS-3716 - Updated Dockerfiles used during integration testing to explicitly install the pip package.  

## [1.18.1] - 2020-07-07

### Fixed
- CASMHMS-3430 - Display disabled nodes correctly with get_xname_status

## [1.18.0] - 2020-06-30

### Added
- CASMHMS-3627 - Updated CAPMC CT smoke test with new API test cases.

## [1.17.6] - 2020-06-09

### Changed
- CASMHMS-3552 - Remove disabled test from CT test deployment.

## [1.17.5] - 2020-05-27

### Changed
- CASMHMS-3502 - Updated readme.md with current implemented features.

## [1.17.4] - 2020-05-05

### Changed
- CASMHMS-2952 - Updated hms-capmc build to use trusted baseOS.

## [1.17.3] - 2020-04-28

### Changed
- Updated the cray-service chart version.

## [1.17.2] - 2020-04-22

### Changed
- CASMHMS-333 - Increase replica count to insure service resiliancy.

## [1.17.1] - 2020-04-02

### Fixed
- CASMHMS-3203 - Fixed regression of SMA cluster name

## [1.17.0] - 2020-03-30

### Changed
- CASMHMS-3211 - Updated get_node_status CT test for disabled nodes.

## [1.16.6] - 2020-03-04

### Fixed
- Updated SMA Postgres cluster address.

# [1.16.5] - 2020-02-22

### Changed
- CASMHMS-2733 - Update dependencies to not use hms-common packages any more.

# [1.16.4] - 2020-02-13

### Fixed
- CASMHMS-2997 - Fixed health endpoint interaction with HSM.

# [1.16.3] - 2020-01-31

### Added
- CASMHMS-2474 - Added initial set of CAPMC Tavern API tests for CT framework.

# [1.16.2] - 2020-01-28

### Fixed
- CASMHMS-2847 - Fixed River telemetry queries to use sql functions we have
  licenses for.

# [1.16.1] - 2020-01-22

### Changed
- CASMHMS-884 - Added liveness, readiness, and health endpoints for the capmc
  service.

# [1.16.0] - 2020-01-07

### Changed
- CASM-1642 - Include River telemetry data in the queries for get_system_power,
  get_system_power_details, get_node_energy, get_node_energy_stats, and
  get_node_energy_counter.

# [1.15.9] - 2019-12-18

### Changed
- CASM-1455 - add routes to 'not implemented' for old Cascade endpoints
  that are not going to be implemented in Shasta

# [1.15.8] - 2019-12-17

### Fixed
- CASMHMS-2576 - get_system_power now returns a better error message
  when there is no data in the query time window
- CASMHMS-2595 - get_node_energy now returns a better error message when
  the start and end times are inverted
- CASMHMS-2612 - get_node_energy_stats now returns a better error message
  when the start and end times are inverted
- CASMHMS-2585 - get_system_power_details now returns a better error message
  when there is no data in the query time window
- CASMHMS-2589 - get_system_power_details now returns a better error message
  when there is no start time specified and no data in the default time window
- CASMHMS-2579 - get_system_power now returns a better error message when there
  is no start time specified and no data in the default time window
- CASMHMS-2598 - get_node_energy_stats now returns a better error message when
  there is no data in the specified time window

# [1.15.7] - 2019-12-16

### Changed
- Updated swagger file to match the code
- Turned off a CT test that needs to be fixed

# [1.15.6] - 2019-12-12

### Changed
- Updated hms-common lib.

# [1.15.5] - 2019-12-11

### Fixed
- CASMHMS-2572 - get_node_energy_counter was returning all the same nid
  values in the results when multiple nids were input.

- CASMHMS-2578 - get_system_power was returning an error code on all
  successful operations.  Now the json return object returns 'e':0
  on success to match the api documentation.

- CASMHMS-2584 - get_system_power_details was returning an error code on 
  all successful operations.  Now the json return object returns 'e':0
  on success to match the api documentation.

- CASMHMS-2577 - get_system_power was incorrectly calculating the default
  start time.  If the user did not pass in a start time, it would query
  a time window that was arbitrarily large.

- CASMHMS-2577 - get_system_power_details was incorrectly calculating the 
  default start time.  If the user did not pass in a start time, it would 
  query a time window that was arbitrarily large.

# [1.15.4] - 2019-12-06

### Added
- CASMHMS-2652 - Add common middleware logging for HTTP incoming requests
  and outbound responses (for the service). There is now a common format
  for both CAPMC as a HTTP client and as a HTTP server logging.

- CASMHMS-2366 - Add new continue request parameter for xname off/on
  APIs which will cause the action to continue with valid component ids
  (xnames) even if the request contains invalid component ids.

- CASMHMS-2653 - Internal refactor of decodeBmcResponse to better
  handle (error) responses from a Redfish API endpoint. The response
  is supposed to be HTTP Content-Type (MIME) application/json but
  there are instances where the body isn't application/json. Use
  the content-type to determine what to do with the response.
  Increases the extesibility of the function to handle other possible
  content-type in the future. Additionally, add support for decoding
  early Dell iDRAC Redfish implementations that don't return a DMTF
  standard Redfish Error (schema) format response.

# [1.15.3] - 2019-11-27

### Fixed
- CASMHMS-2563 - Fixed the DoSystemPower function to not cause a 
  panic when making the query to the databse.  The function now
  also correctly returns an error if there is no data present for
  the requested time interval.

# [1.15.2] - 2019-11-19

### Fixed
- CASMHMS-2246 - Removed the separate Debug variable controlling log
  level information from config.go.  The log level is now determined
  by the global log level.

# [1.15.1] - 2019-11-14

### Changed
- CASMHMS-2172 - Clean up resources and exit cleanly when capmcd
  shuts down.  This change allows the current requests to complete
  and return to the clients before the server stops.  It also
  shuts down the conenctions to resources as cleanly as possible.

# [1.15.0] - 2019-11-12

### Added
- CASMHMS-2471 - Add ability to return Component Id (aka xname) status from
  either Redfish endpoints directly or from the Hardware State Manager. The
  later give a richer set of status but may not necessarily be up to date
  with respect to what Redfish reports. Retrieving state from HSM is much
  faster than going to every Redfish endpoint (controller). The default
  behavior is (still) to go to Redfish.

# [1.14.5] - 2019-11-12

### Changed
- CASMHMS-2327 - Use the HSM discovered Mountain EPO information 
  instead of hard coded paths to perform the EPO. 

# [1.14.4] - 2019-11-11

### Changed
- CASMHMS-2345 - Prefer the PowerControl OData.ID over the PowerControlURL
  returned by HSM for get/set power cap. The later is still the fallback
  when OData.ID isn't present for the control.

# [1.14.3] - 2019-11-08

### Added
- CASMHMS-1545 - Adds a method for specifying the source CAPMC should use
  when determining the status of a node. The pre-CASMHMS-863 behavior can
  be selected by adding `"source": "redfish"` to the JSON payload. The default
  remains using the HSM for status if source is not specified.

# [1.14.2] - 2019-11-06

### Fixed
- CASMHMS-2298 - Return the Cascade NID map format for NIDs not found in
  HSM.  This provides the details on what NIDs were not found.  This now
  returns a format similar to get/set power cap calls.

## [1.14.1] - 2019-11-01

### Fixed
- CASMHMS-2204 - Return the Cascade NID map format for NIDs not found in
  HSM. This provides more detail. The response format is unfortunately not
  documented correctly in the Cascade API document. The Cascade implementation
  does return an extended format (similar to that used by get/set power cap,
  etc.).

## [1.14.0] - 2019-10-31

### Changed
- CASMHMS-863 - The node status (get_node_status) API now uses hardware state
  manager (HSM) component state to determine status. This now matches the
  Cascade get_node_status behavior. CAPMC status is not a one-for-one of the
  defined hardware management system (HMS) states to status, example Disabled
  is not a HMS state but is a CAPMC status. Similarly the Shasta HMS component
  state flags are not a one-for-one match, example Locked is mapped to the
  CAPMC 'resvd' flag.
- Modified CAPMC as a HTTP client logging. It is now possible to only log the
  HTTP body. The HTTP body is logged at debug level 2 or greater. The HTTP
  header is logged at debug level 4 or greater. The body and header are logged
  as expanded strings for even numbered log levels and as Go language safe
  quoted strings for odd numbered log levels.

### Fixed
- Node/Xname status show filter parsing breaks full filter string into tokens
  rather than stopping prematurely. The number of tokens could be greater
  than the number of defined show filters.

## [1.13.1] - 2019-10-22

### Fixed
- CASMHMS-2347 - Ensure get_nid_map returns only enabled and non-empty
  nodes matching Cascade (XC) API behavior.

## [1.13.0] - 2019-10-18

### Added
- CASMHMS-30 - Implement get_power_cap for Mountain hardware
    - Works with systems that support the Redfish Power.PowerControl.PowerLimit
	schema (specifically LimitInWatts)
    - Supports multiple PowerControl (array) items
- CASMHMS-33 - Implement set_power_cap for Mountain hardware
    - Supports proposed Cray Mountain (Windom) nC Redfish Power.PowerControl.PowerLimit via PATCH
    - No support for ETag (Required for Gigabyte)
    - Does not do a read/modify/write (blind PATCH)
    - Uses Cray OEM Redfish extension values retrieved by HSM for value in range checks. If these are not present range validation is skipped.
    - This does not check for valid combination of node and accel because the
      worst that happens is the node will be as slow as a turtle.
- CASMHMS-2283 - Integrate with controller Redfish API GET/PATCH
- Add as a client HTTP request/response logging with additional _middleware_
  logging and new `-debug-level` flag. Logs only when `-debug` is set and
  `-debug-level n` where n > 1. This will log **full** HTTP request/response
  (i.e. includes HTTP headers).
- Unit Test Framework
	- Add ability to skip (`-short`) the 'off' and 'recursive off tests
	  that use a hard coded timer value.
	- Add `-debug` and `-enable-log` flags to unit test framework.
	  These can be used, for example `go test -args -enable-log`, to
	  turn off writing all log.Print* messages to ioutil.Discard.
	  The `-debug` flag sets the equivalent flag that it would in the
	  non-unit test code.

### Changed
- Refactor common code for adding BMC "commands" (API "calls") to the worker
  pool queue. Replace all instances of same loop with new function(s).

## [1.12.0] - 2019-10-09

### Added
- CASMHMS-2102 - Implement get_power_cap_capabilities API/DOMAIN for Shasta.

## [1.11.0] - 2019-10-09

### Added
- CASMHMS-2151 - Enable getting of power cap information from Redfish
- CASMHMS-2152 - Added ability to issue a PATCH to a Redfish endpoint

### Changed
- CASMHMS-2100, CASMHMS-2101 - Validation of get/set power cap APIs
- CASMHMS-2152 - Changed BMC call internals to be more flexible

## [1.10.4] - 2019-10-04

### Fixed

- CASMHMS-2151: Added low level doBmcGetCall function. To allow generic callers
  to query Redfish endpoints for GETs.

## [1.10.3] - 2019-10-03

### Fixed

- CASMHMS-2249: Corrected two printf statements to use %v as required by an interface
- CASMHMS-2248: Updated dockerfile to use golang 1.13 alpine; as just golang:alpine
  apparently does not include 1.13 (which we need to support errors.As syntax)
  Also updated dockerfile to do both coverage and unit tests for ./... to ensure all
  possible unit tests are run

## [1.10.2] - 2019-10-02

### Added

- CASMHMS-29: Added system parameters to config.toml, updated code to allow
  loading of the system parameters

## [1.10.1] - 2019-09-30

### Added

- CASMHMS-1991, CASMHMS-2072, CASMHMS-2076: Added SQL implementation for
  get_node_enegry,   get_node_enegry_stats, get_node_enegry_counter.

## [1.10.0] - 2019-09-27

### Added

- CASMHMS-1689: Component locking for power control commands
- CASMHMS-2174: Add ability to query based on Component/State Enabled flag


## [1.9.3] - 2019-09-24

### Added

- CASMHMS-2074, CASMHMS-2070, CASMHMS-2071: Added API implementation for
  get_node_enegry,   get_node_enegry_stats, get_node_enegry_counter.
- Converted Apid from int to string
- Temporarily put logic in that prevents usage of job_id or apid; as they
  are not supported in PMDB yet

### Fixed

- Fixed the signal handler in main as it didn't exit after closing DB

## [1.9.2] - 2019-09-12

### Fixed

- CASMHMS-897: Add tag for `reason` in power command log messages.

### Changed

- Update database query for `get_system_power_details`. **Only** supports
  Mountain hardware data from the power management database. Removes stub.

### Fixed

- CASMHMS-1937: Respond, via API, with more information from Redfish endpoint
  when commands fail.
- CASMHMS-1898: CAPMC no longer attempts to send Redfish power action reset
  commands to HSM Empty components.

## [1.9.0] - 2019-09-11

### Added

- API to enable Emergency Power Off for Mountain hardware

### Changed

- Update database query for `get_system_power`. **Only** supports Mountain
  hardware data from the power management database. Removes stub.
- Database connection setup, tear-down, and logging rework.

## [1.8.0] - 2019-09-05

### Added

- Add calls for time series database (TSDB) power management database (PMDB)
  - Get total system power (`get_system_power` API) from the PMDB for a given
    start and end time. The values returned are minimum, maximum, and average
    power usage for the given time period.
  - Get total system power per cabinet (`get_system_power_details` API) from
    the PMDB for a given start and end time. The vaules returned are minimum,
    maximum, and average power by cabinet.
  **Note:** These are *stubbed* calls that return the same values regardless of
  input.

## [1.7.5] - 2019-08-27

### Fixed

- Redfish PowerState decoding for `get_xname_status`

## [1.7.4] - 2019-08-29

### Changed

- CAPMC now gets component redfish credentials from Vault.

## [1.7.3] - 2019-08-19

### Fixed

- Use new location of PostgreSQL time series db

## [1.7.2] - 2019-08-14

### Fixed

- Report an error for bad components

## [1.7.1] - 2019-08-08

### Changed

- Update block role list; include Management and no others

## [1.7.0] - 2019-08-07

### Added

- HSM Role based *blacklisting* for power control APIs (off/on/reinit)

## [1.6.0] - 2019-07-23

### Added

- Powering of PDU outlets

## [1.5.0] - 2019-07-05

### Added

- set_power_cap API with a not implemented stub

## [1.4.0] - 2019-07-03

### Added

- get_power_cap API with a not implemented stub

## [1.3.0] - 2019-07-02

### Added

- get_power_cap_capabilities API with a not implemented stub

## [1.2.0] - 2019-05-13

### Changed

- Target v1.1.0 of `hms-common` which no longer includes SMD packages. Updated all those references to now target SMD repo.

## [1.1.0] - 2019-05-13

### Fixed

- Version file now tracked.

## [1.0.0] - 2019-05-13

### Added

- This is the initial release. It contains everything that was in `hms-services` at the time with the major exception of being `go mod` based now.
