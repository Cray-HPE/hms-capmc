1.  [CASMHMS](index.html)
2.  [CASMHMS Home](CASMHMS-Home_119901124.html)
3.  [Design Documents](Design-Documents_127906417.html)
4.  [CAPMC](CAPMC_227479290.html)

# <span id="title-text"> CASMHMS : CAPMC Power Documentation </span>

Created by <span class="author"> Michael Jendrysik</span>, last modified
on Dec 11, 2020

-   [Components capable of having power
    controlled](#CAPMCPowerDocumentation-Componentscapableofhavingpowercontrolled)
    -   [River](#CAPMCPowerDocumentation-River)
    -   [Mountain](#CAPMCPowerDocumentation-Mountain)
-   [Groupings](#CAPMCPowerDocumentation-Groupings)
    -   [Component groups](#CAPMCPowerDocumentation-Componentgroups)
    -   [Specials and high level
        components](#CAPMCPowerDocumentation-Specialsandhighlevelcomponents)
    -   [Wild cards](#CAPMCPowerDocumentation-Wildcards)
-   [Power sequencing](#CAPMCPowerDocumentation-Powersequencing)
-   [Talking to CAPMC](#CAPMCPowerDocumentation-TalkingtoCAPMC)
    -   [Cray CLI](#CAPMCPowerDocumentation-CrayCLI)
        -   [For example, to power off a node via
            xname](#CAPMCPowerDocumentation-Forexample,topoweroffanodeviaxname)
        -   [To power on a Chassis and all components below
            it](#CAPMCPowerDocumentation-TopoweronaChassisandallcomponentsbelowit)
    -   [curl](#CAPMCPowerDocumentation-curl)
        -   [For example, to power off a node via
            xname](#CAPMCPowerDocumentation-Forexample,topoweroffanodeviaxname.1)
        -   [To power on a Chassis and all components below
            it](#CAPMCPowerDocumentation-TopoweronaChassisandallcomponentsbelowit.1)
-   [Power capping Admin guide
    updates](#CAPMCPowerDocumentation-PowercappingAdminguideupdates)
    -   [Section 17.3 of v1.3
        RevE](#CAPMCPowerDocumentation-Section17.3ofv1.3RevE)

<span class="jira-issue resolved" jira-key="CASMHMS-1339">
<a href="https://connect.us.cray.com/jira/browse/CASMHMS-1339?src=confmacro" class="jira-issue-key"><img src="https://connect.us.cray.com/jira/secure/viewavatar?size=xsmall&amp;avatarId=13315&amp;avatarType=issuetype" class="icon" />CASMHMS-1339</a>
- <span class="summary">Create Customer Facing Documentation of what
gets powered off </span> <span
class="aui-lozenge aui-lozenge-subtle aui-lozenge-success jira-macro-single-issue-export-pdf">Done</span>
</span>

<span class="jira-issue resolved" jira-key="CASMHMS-1340">
<a href="https://connect.us.cray.com/jira/browse/CASMHMS-1340?src=confmacro" class="jira-issue-key"><img src="https://connect.us.cray.com/jira/secure/viewavatar?size=xsmall&amp;avatarId=13315&amp;avatarType=issuetype" class="icon" />CASMHMS-1340</a>
- <span class="summary">Create documentation of what will be powered
on/off via CAPMC</span> <span
class="aui-lozenge aui-lozenge-subtle aui-lozenge-success jira-macro-single-issue-export-pdf">Done</span>
</span>

# <span style="color: rgb(0,0,0);">Components capable of having power controlled</span>

## River

**Compute Nodes**: x0c0s0b0n0

**iPDU sockets**: x0p0j0 - POR but not currently supported in CAPMC,
support coming in Q2'19

**non-compute nodes**: x0c0s0b0n0 - POR but not currently supported in
CAPMC, unknown when support is coming

## Mountain

**Chassis:** x0c0

**Switch Modules**: x0c0r0

**Compute Modules**: x0c0s0

**Nodes**: x0c0s0b0n0

# Groupings

### Component groups

**Component groups:**

The CMS team plans to use node groups. The group information would be
stored in the HSM and be accessible to the other running services. CAPMC
currently supports on/off by group and reinit by group is coming soon.

### Specials and high level components

**Full system**: s0, all

**Cabinets**: x0, x1000, x2500, etc...

**Chassis**: x0c0, x1000c0, x2500c7, etc...

**Modules**: x0c0s0, x1000c0s3, x2500c7s7, etc...

By default, CAPMC handles power on 1 and only 1 component at a time.
There is an option called ***recursive*** that ca be passed to CAPMC via
craycli or the curl interface. Please
see <a href="http://web.us.cray.com/~pubs/cray-portal/public/iaas/capmc/overview/" class="external-link">Cray Advanced Platform Monitoring and Control (CAPMC) API Overview</a>
for more details on the RESTful interface. Under the covers craycli
makes calls to the RESTful interface. When the recursive option is
included in a request, all of the sub-components of the target component
are included in the power command.

Ancestors option

The ancestors option **prereq** will allow the admin to tell CAPMC to
power on a node (x0c0s0b0n0) and CAPMC will make sure all of the
components above that node are also powered on. So on a cold cabinet,
the target list inside of CAPMC changes from x0c0s0b0n0 to x0c0, x0c0s0,
x0c0s0b0n0. We have not had a discussion on if the switch modules should
also be powered on, that is a future enhancement at this time.

### Wild cards

**Wild carding**: POR but not currently supported in CAPMC, unknown when
support is coming

x\[0-1\] - cabinets 1 and 2

x0c\[0, 2,4,6\] - all the chassis one (left?) side of a cabinet

x\[0-3\]c\[0-3\]s\[0-7\] - all slots in the top 4 chassis of cabinets 0,
1, 2, and 3

# Power sequencing

CAPMC assumes the all cabinets and iPDUs have been plugged in, breakers
are on, and iPDU controller and CMMs for all Chassis are on and
available.

CAPMC contains a default ordering of components for powering on. This
ordering will eventually be configurable. Support for that will be in by
the end of Q2'19

Power ON Mountain:Chassis, Switch Modules, Compute Modules, Nodes

Power OFF Mountain: Nodes, Compute Modules, Switch Modules, Chassis

Power ON River: iPDU Outlet, Nodes

Power OFF River: Nodes, iPDU Outlet

When a mix of xnames are sent to CAPMC or one of the special names or
groups is used, the components will be ordered based on their type. The
target components will then be powered on/off based on the ordering
above.

  

# Talking to CAPMC

## Cray CLI

The Cloud team has developed a CLI command that is based loosely on
other cloud commands such as AWS and gcloud. Below are some of the help
outputs for capmc. There is new code we are trying to get in today that
I don't have output for, but I did include an example.

The Cray CLI is designed to be executed from any system that has https
access to the SMS of the system. There is a gateway service that runs on
the SMS and provides the means to contact the other services that are
running on the SMS such as CAPMC. Please contact the Cloud team for more
information on how to configure the Cray CLI, which systems it can
execute on, and authentication information.

    > cray
    Usage: cray [OPTIONS] COMMAND [ARGS]...

      Cray management and workflow tool

    Options:
      --version  Show the version and exit.
      --help     Show this message and exit.

    Groups:
      ars       Artifact Repository Service API
      auth      Manage OAuth2 credentials for the Cray CLI
      badger    Badger
      bss       Boot Script Service API
      capmc     Cray Advanced Platform Monitoring and Control (CAPMC) API
      cfs       Configuration Framework Service
      config    View and edit Cray configuration properties
      fabric    Fabric Manager REST API
      firmware  Firmware Update Service
      hsm       Hardware State Manager API
      ims       Image Management Service (IMS) API
      nmd       Cray Ckdump Service
      pals      Parallel Application Launch Service API
      prs       Cray Package Repository Service
      reds      Shasta River Node Initialization, Geolocation and RedFish...
      sds       framework required base group
      uas       Cray User Access Service

    Commands:
      init     Initialize/reinitilize the Cray CLI
      mpiexec  Run an application using the Parallel Application Launch Service

    > cray capmc
    Usage: cray capmc [OPTIONS] COMMAND [ARGS]...

      Cray Advanced Platform Monitoring and Control (CAPMC) API

    Options:
      --help  Show this message and exit.

    Groups:
      get_nid_map
      get_node_rules
      get_node_status
      get_xname_status
      node_off
      node_on
      node_reinit
      xname_off
      xname_on


    > cray capmc xname_off
    Usage: cray capmc xname_off [OPTIONS] COMMAND [ARGS]...

    Options:
      --help  Show this message and exit.

    Commands:
      create


    > cray capmc xname_off create --help
    Usage: cray capmc xname_off create [OPTIONS]

    Options:
      --force BOOLEAN            This flag disables any checks for a graceful
                                 power off. The action is taken immediately.
      --xnames TEXT,...          User specified list of xnames to shutdown and
                                 power off. An empty array is invalid. If invalid
                                 xnames are specified then an error will be
                                 returned.  [required]
      --reason TEXT              Arbitrary, free-form text
      --configuration CONFIG     name of configuration to use. Create through
                                 `cray init`  [env var: CRAY_CONFIG; required]
      --quiet
      --format [json|toml|yaml]
      --token TOKEN_FILE_PATH    [env var: CRAY_CREDENTIALS]
      -v, --verbose              Example: -vvvv
      --help                     Show this message and exit.

### For example, to power off a node via xname

    cray capmc xname_off create --xnames x0c0s0b0n0

Powers off the node x0c0s0b0n0

### To power on a Chassis and all components below it

    cray capmc xname_on create --xnames x0c0 --recursive

Powers on the chassis x0c0, all the switch modules x0c0r\[0-7\], all the
compute modules x0c0s\[0-7\], and all the nodes
x0c0s\[0-7\]b\[0-1\]n\[0-1\]

## curl

The other way to talk to CAPMC is via curl. This can be done remotely
through the gateway or from the SMS via the gateway. There are other
ways, but due to the authentication being put in place you will need to
talk to the gateway. Also, since we do not currently have authentication
in place, I do not know what extra commands are needed to get the proper
credentials to talk to CAPMC via curl.

Here are the curl versions of the above examples

### For example, to power off a node via xname

    curl -k -X POST https://slice-sms:30443/apis/capmc/capmc/xname_off -d '{"xnames":["x0c0s0b0n0"]}'

Powers off the node x0c0s0b0n0

### To power on a Chassis and all components below it

    curl -k -X POST https://slice-sms:30443/apis/capmc/capmc/xname_on -d '{"xnames":["x0c0"], "recursive":true}'

Powers on the chassis x0c0, all the switch modules x0c0r\[0-7\], all the
compute modules x0c0s\[0-7\], and all the nodes
x0c0s\[0-7\]b\[0-1\]n\[0-1\]

# Power capping Admin guide updates

## Section 17.3 of v1.3 RevE

**17.3 Air Cooled Node Power Management**

Air Cooled node power management is supported by the server BMC
firmware. The BMC exposes the power control API for a node via the
node's Redfish ChassisPower schema. Out-of-band power management data is
polled by a collector and published on a kafka bus for entry into the
Power Management Database.

The Cray Advanced Platform Management and Control (CAPMC) API
facilitates power control and enables power aware WLMs such as Slurm to
perform power management and power capping tasks.

**Always use the Boot Orchestration Service (BOS) to power off or power
on compute nodes.**

Redfish API

The Redfish API for Air Cooled nodes is the node's Chassis Power
resource which is presented by the BMC. OEM properties may be used to
augment the Power schema and allow for feature parity with previous Cray
system power management capabilities. A PowerControl resource presents
the various power management capabilities for the node.

Each node has one power control resource

-   Node power control

The power control of the node will need to be enabled an may require
additional licenses to use.

(@jmo Not sure what to do here. CAPMC is not responsible for enabling
power capping on nodes. A different process for each server vendor will
be needed. Your Active and Deactivate below are for Gigabyte nodes
only.)

**Power Capping**

CAPMC power capping controls for compute nodes can query component
capabilities and manipulate the node power constraints. This
functionality enables external software to establish an upper bound, or
estimate a minimum bound, on the amount of power a system may consume.

CAPMC API calls provide means for third party software to implement
advanced power management strategies and JSON functionality can send and
receive customized JSON data structures.

Air Cooled nodes support these power capping and monitoring API calls:

-   get\_power\_cap\_capabilities
-   get\_power\_cap
-   set\_power\_cap
-   get\_node\_energy
-   get\_node\_energy\_stats
-   get\_system\_power

**Cray CLI Examples for Liquid Cooled Compute Node Power Management**

Get Node Energy

    ncn-w001# cray capmc get_node_energy create --nids NID_LIST --start-time '2020-03-04 12:00:00' --end-time '2020-03-04 12:10:00' --format json

Get Node Energy Stats

    ncn-w001# cray capmc get_node_energy_stats create --nids NID_LIST --start-time '2020-03-04 12:00:00' --end-time '2020-03-04 12:10:00' --format json

Get Node Power Control and Limit Settings

    ncn-w001# cray capmc get_power_cap create –-nids NID_LIST --format json

Get System Power

    ncn-w001# cray capmc get_system_power create --start-time '2020-03-04 12:00:00' --window-len 30 --format json

Get Power Capping Capabilities

The supply field contains the Max limit for the node.

    ncn-w001# cray capmc get_power_cap_capabilities create –-nids NID_LIST --format json

Set Node Power Limit

    ncn-w001# cray capmc set_power_cap create –-nids NID_LIST --node 225 --format json

Remove Node Power Limit (Set to Default)

    ncn-w001# cray capmc set_power_cap create –-nids NID_LIST --node 0 --format json

  

Activate Node Power Limit

\# curl -k -u $login:$pass -H "Content-Type: application/json" \\  
-X POST
https://$BMC\_IP/redfish/v1/Chassis/Self/Power/Actions/LimitTrigger
--date  
'{"PowerLimitTrigger": "Activate"}'

Deactivate Node Power Limit

\# curl -k -u $login:$pass -H "Content-Type: application/json" \\  
-X POST
https://$BMC\_IP/redfish/v1/Chassis/Self/Power/Actions/LimitTrigger
--data  
'{"PowerLimitTrigger": "Deactivate"}'

Document generated by Confluence on Jan 14, 2022 07:17

[Atlassian](http://www.atlassian.com/)
