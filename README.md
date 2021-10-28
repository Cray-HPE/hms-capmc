# Cray Advanced Platform Monitoring and Control (CAPMC) Service

Cray Advanced Platform Monitoring and Control provides a way to monitor and
control certain components in a Shasta system. CAPMC uses a RESTful interface to
provide monitoring and control capabilities and executes in the management plane
in the SMS cluster. Administrator level permissions are required for most
operations. The CAPMC service relies on a running Hardware State Manager (HSM)
service. The HSM contains all of the necessary information for CAPMC to
communicate with the hardware.

## Building and Executing CAPMC

### Building CAPMC

[Building CAPMC after the Repo split](https://connect.us.cray.com/confluence/display/CASMHMS/HMS+Repo+Split)

### Running CAPMC (capmcd) locally

Starting capmcd:

```bash
./capmcd -http-listen="localhost:27777" -hsm=https://localhost:27779
```

#### Caveats: Connecting to the TSDB
By default the DB connection will try to connect to Postgres.  Use the following ENV VARs to specify where to  try to connect:
```
DB_HOSTNAME=somePostgresDB
DB_PORT=thePort
```

Example of CURL command to make sure it is working:

```bash
curl -X POST -i -d '{"nids":[7]}' http://localhost:27777/capmcd/get_node_status
```

### Running CAPMCD in Docker Container

From the root of this repo, build the image:

```bash
docker build -t cray/capmcd:1.0 .
```  

Then run (add `-d` to the arguments list of `docker run` to run in detached/background mode):

```bash
docker run -p 27777:27777 --name capmcd cray/capmcd:1.0
```

All connections to localhost on port 27777 will flow through the running container.

### Using CAPMC through the Cray CLI interface

Example to power on an entire cabinet:

```bash
cray capmc xname_on create --xnames x1000
```

Example to power off a Chassis an all of its descendents:

```bash
cray capmc xname_off create --xnames x1000c0 --recursive
```

## Build, Tag, Push

./build_tag_push.sh -l <host system>:5000

On target system, delete the running pod and the one pushed will get started.

## CAPMC CT Testing

This repository builds and publishes hms-capmc-ct-test RPMs along with the service itself containing tests that verify CAPMC on
the NCNs of live Shasta systems. The tests require the hms-ct-test-base RPM to also be installed on the NCNs in order to execute.
The version of the test RPM installed on the NCNs should always match the version of CAPMC deployed on the system.

## API Map

When the different APIs will be supported:

| Equivalent XC | v1 now | v1 future |
| --- | --- | --- | --- |
| get_nid_map | get_nid_map | - |
| get_node_rules | get_node_rules | - |
| get_node_status | get_node_status | - |
| node_on | node_on | - |
| node_off | node_off | - |
| node_reinit | node_reinit | - |
| - | get_xname_status | - |
| - | xname_on | - |
| - | xname_off | - |
| - | xname_reinit | - |
| - | group_on | - |
| - | group_off | - |
| - | group_reinit | - |
| - | get_group_status | - |
| - | emergency_power_off | - |
| get_power_cap_capabilities | get_power_cap_capabilities | - |
| get_power_cap | get_power_cap | - |
| set_power_cap | set_power_cap | - |
| get_node_energy | get_node_energy | - |
| get_node_energy_stats | get_node_energy_stats | - |
| get_node_energy_counter | get_node_energy_counter | - |
| get_system_power | get_system_power | - |
| get_system_power_details | get_system_power_details | - |
| get_system_parameters | get_system_parameters | - |
| get_partition_map | - | get_partition_map |
| - | - | get_partition_status |
| - | - | partition_on |
| - | - | partition_off |
| - | - | partition_reinit |
| - | - | get_gpu_power_cap_capabilities |
| - | - | get_gpu_power_cap |
| - | - | set_gpu_power_cap |
| get_power_bias | - | get_power_bias (if needed) |
| set_power_bias | - | set_power_bias (if needed) |
| clr_power_bias | - | clr_power_bias (if needed) |
| set_power_bias_data | - | set_power_bias_data (if needed) |
| compute_power_bias | - | compute_power_bias (if needed) |
| get_freq_capabilities | - | get_freq_capabilities (if needed ) |
| get_freq_limits | - | get_freq_limits (if needed) |
| set_freq_limits | - | set_freq_limits (if needed) |
| get_sleep_state_limite_capabilities | - | get_sleep_state_limite_capabilities (if needed) |
| set_sleep_state_limit | - | set_sleep_state_limit (if needed) |
| get_sleep_state_limit | - | get_sleep_state_limit (if needed) |
| get_mcdram_capabilities (Xeon Phi) | - | - |
| get_mcdram_cfg (Xeon Phi) | - | - |
| set_mcdram_cfg (Xeon Phi) | - | - |
| clr_mcdram_cfg (Xeon Phi) | - | - |
| get_numa_capabilities (Xeon Phi) | - | - |
| get_numa_cfg (Xeon Phi) | - | - |
| set_numa_cfg (Xeon Phi) | - | - |
| clr_numa_cfg (Xeon Phi) | - | - |
| get_ssd_enable (XC Only) | - | - |
| set_ssd_enable (XC Only) | - | - |
| clr_ssd_enable (XC Only) | - | - |
| get_ssds (XC Only) | - | - |
| get_ssd_diags (XC Only) | - | - |

## Current Features

* Power control
  * Redfish power status of components
  * Single components via NID or xname
  * Grouped components
  * Entire system (all or s0)
  * Per cabinet (x1000)
  * Ancestors and descendants of single component
  * Force option for immediate power off
  * Node power capping
  * Emergency Power Off at the Chassis level
  * Query of power data at node, system, and cabinet level

## Future Features and updates

* Backend performance improvements
* Moving to a truly RESTful interface (v2)
* Power control
  * Emergency Power Off at the iPDU levels
  * Power control of Mountain CDUs (won't/cant do)
  * Power control policies
  * Power control of Motivair door fans
  * Power control of in-rack River CDUs
* Power capping and related for Mountain
  * Group level and system level power capping (if needed)
  * Power bias factors to individual nodes (if needed)
  * Query of power data at group level (if needed)
  * RAPL (Running Average Power Limiting) (if possible)
* Node level CState/Pstate handling (if needed and not handled by WLM)
* GPU power capping
* Powering off idle nodes (most likely a WLM function)
* Rebooting nodes (most likely a CMS or WLM function)

## Limitations

* No Redfish interface to control Mountain CDUs
* CMM and CEC cannot be powered off. They are always ON when Mountain cabinets
  are plugged in and breakers are ON
* Can only talk to components that exist in HSM
