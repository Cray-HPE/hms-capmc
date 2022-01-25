1.  [CASMHMS](index.html)
2.  [CASMHMS Home](CASMHMS-Home_119901124.html)
3.  [Design Documents](Design-Documents_127906417.html)
4.  [CAPMC](CAPMC_227479290.html)

# <span id="title-text"> CASMHMS : Shasta v1.4 Power Management Guide </span>

Created by <span class="author"> Michael Jendrysik</span>, last modified
on Apr 16, 2021

<table class="wrapped confluenceTable">
<tbody>
<tr class="odd">
<th class="confluenceTh">Author</th>
<td class="confluenceTd"><div class="content-wrapper">
<p><a href="https://connect.us.cray.com/confluence/display/~mjendrysik" class="confluence-userlink user-mention">Michael Jendrysik</a></p>
</div></td>
</tr>
<tr class="even">
<th class="confluenceTh">Publications</th>
<td class="confluenceTd"></td>
</tr>
<tr class="odd">
<th class="confluenceTh">Product</th>
<td class="confluenceTd"><em>Shasta</em></td>
</tr>
<tr class="even">
<th class="confluenceTh">Related</th>
<td class="confluenceTd"><div class="content-wrapper">
<p><span class="jira-issue resolved" data-jira-key="CASMHMS-4583"> <a href="https://connect.us.cray.com/jira/browse/CASMHMS-4583?src=confmacro" class="jira-issue-key"><img src="https://connect.us.cray.com/jira/secure/viewavatar?size=xsmall&amp;avatarId=13315&amp;avatarType=issuetype" class="icon" />CASMHMS-4583</a> - <span class="summary">Provide WhitePaper on power Management Capabilities</span> <span class="aui-lozenge aui-lozenge-subtle aui-lozenge-success jira-macro-single-issue-export-pdf">Done</span> </span></p>
<p><a href="https://support.hpe.com/hpesc/public/docDisplay?docId=a00111788en_us" class="external-link">https://support.hpe.com/hpesc/public/docDisplay?docId=a00111788en_us</a> (local copy below)</p>
<p><a href="https://connect.us.cray.com/jira/secure/attachment/102893/HPE_a00111788en_us_HPE%20Performance%20Cluster%20Manager%20Power%20Management%20Guide.pdf" class="external-link">https://connect.us.cray.com/jira/secure/attachment/102893/HPE_a00111788en_us_HPE%20Performance%20Cluster%20Manager%20Power%20Management%20Guide.pdf</a></p>
</div></td>
</tr>
<tr class="odd">
<th class="confluenceTh">Target release</th>
<td class="confluenceTd">v1.4</td>
</tr>
<tr class="even">
<th class="confluenceTh">Date</th>
<td class="confluenceTd"><em>2021/03/30</em></td>
</tr>
<tr class="odd">
<th class="confluenceTh">Document status</th>
<td class="confluenceTd"><div class="content-wrapper">
<span class="status-macro aui-lozenge">DRAFT</span>
</div></td>
</tr>
</tbody>
</table>

# <span class="nh-number">1. </span>Notices

The information contained herein is subject to change without notice.
The only warranties for Hewlett Packard Enterprise products and services
are set forth in the express warranty statements accompanying such
products and services. Nothing herein should be construed as
constituting an additional warranty. Hewlett Packard Enterprise shall
not be liable for technical or editorial errors or omissions contained
herein.

C<span class="inline-comment-marker"
ref="0e5ac3ab-f163-43f6-ac7a-3e294efba4f7">onfidential computer
software. Valid license from Hewlett Packard Enterprise required for
possession, use, or copying. Consistent with FAR 12.211 and 12.212,
Commercial Computer Software, Computer Software Documentation, and
Technical Data for Commercial Items are licensed to the U.S. Government
under vendor's standard commercial license.</span>

Links to third-party websites take you outside the Hewlett Packard
Enterprise website. Hewlett Packard Enterprise has no control over and
is not responsible for information outside the Hewlett Packard
Enterprise website.

# <span class="nh-number">2. </span>Acknowledgements

# <span class="nh-number">3. </span>Change History

<table class="wrapped confluenceTable">
<tbody>
<tr class="odd">
<td class="confluenceTd"><p><strong>Date</strong></p></td>
<td class="confluenceTd"><p><strong>Version</strong></p></td>
<td class="confluenceTd"><p><strong>Change Description</strong></p></td>
</tr>
<tr class="even">
<td class="confluenceTd"><p><em>2021/04/13</em></p></td>
<td class="confluenceTd"><p><em>0.0</em></p></td>
<td class="confluenceTd"><p>Initial version, supports the Shasta 1.4 release</p></td>
</tr>
</tbody>
</table>

**Table of Contents**

-   [1. Notices](#Shastav1.4PowerManagementGuide-Notices)
-   [2.
    Acknowledgements](#Shastav1.4PowerManagementGuide-Acknowledgements)
-   [3. Change History](#Shastav1.4PowerManagementGuide-ChangeHistory)
-   [4. Shasta Power
    Management](#Shastav1.4PowerManagementGuide-ShastaPowerManagement)
    -   [4.1. Abstract](#Shastav1.4PowerManagementGuide-Abstract)
    -   [4.2. Overview](#Shastav1.4PowerManagementGuide-Overview)
    -   [4.3. Power management
        features](#Shastav1.4PowerManagementGuide-Powermanagementfeatures)
        -   [4.3.1. Supported
            Hardware](#Shastav1.4PowerManagementGuide-SupportedHardware)
    -   [4.4. Power Management
        Controls](#Shastav1.4PowerManagementGuide-PowerManagementControls)
        -   [4.4.1. Component
            Control](#Shastav1.4PowerManagementGuide-ComponentControl)
        -   [4.4.2. Power
            Capping](#Shastav1.4PowerManagementGuide-PowerCapping)
        -   [4.4.3. Node
            Energy](#Shastav1.4PowerManagementGuide-NodeEnergy)
        -   [4.4.4. System
            Power](#Shastav1.4PowerManagementGuide-SystemPower)
        -   [4.4.5. Emergency Power
            Off](#Shastav1.4PowerManagementGuide-EmergencyPowerOff)
-   [5. Future Shasta Power
    Management](#Shastav1.4PowerManagementGuide-FutureShastaPowerManagement)

 

# <span class="nh-number">4. </span><u>Shasta Power Management </u>

## <span class="nh-number">4.1. </span>Abstract

Cray System Management (CSM) provides microservices for power control
and "power capping" of a Cray EX Supercomputer. This paper aims to
clarify what power control is available at the component, node, chassis,
cabinet, and system levels. Administrators familiar with other HPC
systems that support discrete power control will find parallels in the
CSM microservices and should be able to use this paper to plan and
manage power consumption once the Cray EX Supercomputer is installed on
site.

The main power management API in CSM 1.4 is called CAPMC and is heavily
inspired by a legacy Cray application from the Cascade class systems. In
CSM 1.5, CAPMC will be supplanted by a new API with more deterministic
behavior, but the overall feature set will be retained. In this paper,
CAPMC terminology will be used to describe the feature set and a
subsequent whitepaper will be released with CSM 1.5 to update applicable
terminology.

In keeping with the Shasta vision, all hardware actions use Redfish APIs
and where the device fails to support the Redfish standard, a translator
is used to allow CSM software to focus on a common interface. As such,
all CAPMC commands related to hardware are ultimately expressed in
Redfish. Even with translators in place, different devices support
different levels of power control which limits the degree to which power
management can be used on each device.

## <span class="nh-number">4.2. </span>Overview

Shasta power management is based on the Cascade power management feature
set. The ability to turn components on and off, limiting power of a node
and/or accelerator, and <span class="inline-comment-marker"
ref="651ff2fa-f328-492b-8f1e-cdf765c2712c">query historical energy and
system power information are all supported in Shasta</span>. The Cray
Advanced Platform Monitoring and Control (CAPMC) service provides all of
these capabilities. An administrator uses the CrayCLI to interact with
CAPMC and third-party software, such as workload managers, can go
directly to the CAPMC API. These controls provide the ability for an
administrator or workload manager (WLM) to manage the power state of an
entire system and to more intelligently manage systemwide power
consumption. As systems grow in size, they grow in power draw as well.
This power draw can be a burden on the local power grid during node
boot, job start, job stop, and node shutdown. By limiting compute node
power across a system, a system wide power limit can be achieved. If
there are compute nodes that have been idle long enough, they can be
powered down by the WLM to help conserve power.

## <span class="nh-number">4.3. </span>Power management features

CAPMC runs as a service inside the Kubernetes service mesh on the worker
NCNs to interact with the various systems and their power control
interfaces. These power control interfaces that CAPMC interacts with
speak the Redfish protocol using a RESTful set of APIs. The CAPMC to
Redfish communication is generally hidden from the caller except in the
case of an error. There are three instances of the service running at
one time and each CAPMC request will be handled completely by one of
those three instances. CAPMC relies heavily on the contents of the
Hardware State Manager (HSM). When CAPMC receives a request it queries
the HSM to determine the hostname of the controllers it needs and what
URLs on those controllers it needs to complete the callers request.
After consulting with HSM, CAPMC will contact each controller to
complete the the request.

<span class="nolink">API documentation in docker form
(shastadocs-website-1.4.0-20210312162254\_c2b0a47-dockerimage.tar) and
directory form
(shastadocs-directory-1.4.0-20210312162254\_c2b0a47.tar.gz) are
available and packaged with the Shasta product RPMs.  
</span>

The CAPMC service performs the following functions:

-   Powers components on and off by single component, NID, and group
-   Power status
-   Node level power capping
-   <span class="inline-comment-marker"
    ref="c3dfbf23-e086-40d3-9b8c-9e656706822c">NID to xname
    mapping</span>
-   <span class="inline-comment-marker"
    ref="9f63cb01-7662-4ddc-ae10-19999711bef8">Node level energy</span>
-   <span class="inline-comment-marker"
    ref="9f63cb01-7662-4ddc-ae10-19999711bef8">System level power
    information</span>
-   <span class="inline-comment-marker"
    ref="9f63cb01-7662-4ddc-ae10-19999711bef8">Per cabinet power
    information</span>
-   Emergency Power Off

### <span class="nh-number">4.3.1. </span>Supported Hardware

Devices support power control through their BMCs which CAPMC can contact
via the exposed Redfish APIs. In Shasta 1.4, several classes of devices
have been tested and are supported at various levels.

The following hardware is supported by CAPMC for power On/Off/Status
commands:

-   Cray EX liquid-cooled hardware
-   HPE ProLiant DL325/385 management servers
-   Gigabyte compute servers
-   SeverTech iPDU **(<span
    class="status-macro aui-lozenge aui-lozenge-current">TODO</span>NEED
    VERSIONS WE SHIP WITH)**

The following hardware is supported by CAPMC for Emergency Power Off:

-   Cray EX liquid-cooled hardware

The following systems support power management at the node level:

-   Cray EX Windom compute nodes
-   Gigabyte compute servers

The following systems support power management of accelerators (GPUs):

-   Cray EX Windom compute nodes

## <span class="nh-number">4.4. </span>Power Management Controls

### <span class="nh-number">4.4.1. </span>Component Control

One of the main operations of CAPMC is to handle power control of
different components. These controls enable powering Off or On
components (chassis, modules, nodes, etc...), and querying component
state information using component identifiers known as *xnames,* node
IDs (NID), or group names. Power control of components is only available
to those components that have been discovered by HSM.

Only certain types of hardware is controllable by CAPMC and can have its
state queried:

-   Cray EX liquid-cooled hardware: Nodes, Node slots, HSN boards,
    Switch slots, and Chassis
-   Cray EX air-cooled hardware: Nodes, HSN Boards, and iPDU Outlets

When querying node power state by NID, CAPMC will query the Hardware
State Manager for its knowledge of the node state. This is done to
facilitate a Workload Managers use of the power state information. All
Hardware State Manager states can be returned including **Ready**,
**Standby**, **On**, and **Off**. In addition, filters can be used to
limit the response.

When querying node power state by xname or group name, CAPMC will query
the hardware for its current state. The only power states that will be
returned are **On**, and **Off**. In addition, filters can be used to
limit the response.

The behavior for all Redfish actions is as follows:

-   Query HSM for BMC connection information about the target
    components.
-   Query Vault for Redfish credentials.
-   Initiate an HTTPS request to the BMCs of the target components to
    perform the operation.

The most common components a power operation would be enacted on are:
Nodes, Node slots, Switch slots, and Chassis. It is rare that direct
power control of a HSN board or iPDU Outlet is needed. All CAPMC power
operations should be viewed as asynchronous and require the client to
check the status of the hardware after the initial power command
returns. After sending the talking to the BMC, CAPMC will return with a
success that the BMC request succeeded or an error indicating what went
wrong. Even if the CAPMC request to the BMC succeeded, that does not
guarantee the power state of the target component will change. Hardware
problems and external influencers can prevent a component from properly
powering On or Off.

Certain components require a higher level component to be powered on
first. In the event that is the case, CAPMC will return an error
indicating it could not talk to the BMC. This error may look similar to
a network connectivity error.

Notes:

-   Powering off a liquid-cooled node slot will drop the power to all
    BMCs and nodes in that slot without a graceful shutdown. There is a
    method to shut down components recursively.
-   The Chassis Management Module (CMM) protects against powering off a
    Chassis if there are switch slots or node slots powered on. A
    force-off command overrides this feature.
-   Powering off the wrong iPDU Outlets can prevent access to an
    air-cooled rack if those outlets are powering the TOR Ethernet
    switch.
-   Emergency Power Off is for Cray EX liquid-cooled hardware only.
-   In an Cray EX liquid-cooled hardware; Chassis, switch slots, and
    node slots are controlled by the Chassis BMC.
-   All nodes are controlled by a node BMC.
-   The ChassisBMC is part of the CMM that has its power controlled by
    breakers on the Cray EX liquid-cooled cabinets.

Additionally, a mechanism is provided for a system administrator to
convey information about hardware (and perhaps site-specific) rules or
timing constraints that allow for efficient and effective management of
idle node resources. The data returned informs the client of how long it
is expected to take to power On or Off a node, the minimum amount of
time nodes should be left off to save energy, and limits on the number
of nodes that should be turned on or off at once. Default rules are
supplied where appropriate. Other values such as the maximum node counts
for power On and Off operations and the maximum amount of time a node
should remain off after a power down are left unset. The values are not
strictly enforced by CAPMC as they are meant to provide guidelines for
authorized clients in their use of the CAPMC service.

### <span class="nh-number">4.4.2. </span>Power Capping

Power capping controls enable querying component capabilities and
manipulation of node or sub-node (accelerator) power constraints. This
functionality enables external software to establish an upper bound, or
estimate a minimum bound, on the amount of power a system or a select
subset of the system may consume. Power cap capabilities are available
on a per-node type basis providing min and max values allowed for power
capping of both a node and any accelerators that are installed. NCNs may
not be power capped and any attempt to do so will result in an error.
System level power limiting is an administrative task by applying
node-level power capping to selected compute nodes.

Cray EX liquid-cooled compute nodes have power capping enabled by
default and set to *unlimited* (no power cap). Power capping for HPE
Proliant DL and Gigabyte servers requires an additional license. Please
refer to the iLO and/or the MegaRack documentation on how to acquire and
apply a license and enable power capping for those systems.

### <span class="nh-number">4.4.3. </span><span class="inline-comment-marker" ref="fda743d4-78e8-4d43-8103-09167940af4e">Node Energy</span>

To help understand power usage of a system, the accumulated node energy
consumption of a set of nodes may be queried by a time range. Types of
information returned may include aggregated energy usage on a set of
nodes, energy usage per-node, or an energy accumulation point in time
snapshot.

Node energy calls backed by data from the power management database
(PMDB), are subject to a hysteresis constraint. The hysteresis value is
15 seconds; so all data sets returned from those API calls will be
shifted back in time so that the end time is less than the hysteresis
window. The minimum window of time that may be queried is 15 seconds.
Furthermore, all aggregation of data over the time dimension will first
be aggregated in 15s time windows. These constraints ensure high data
quality and minimize the impact of hysteresis on API calls.

Node energy queries can be resource intensive. Depending on system size
and input parameters, those calls may require several minutes to
complete.

### <span class="nh-number">4.4.4. </span>System Power

CAPMC provides a way to monitor near real time system level power
consumption. Current or historical power data may be selected. This data
may be returned in aggregate or constituent form containing information
relating to total system or per-cabinet components, respectively.
Additionally, a mechanism is provided for a system administrator to
convey intent or other operational parameters, such as a maximum system
power limit, or unreported static power overhead. The values are not
strictly enforced by CAPMC as they are meant to provide guidelines for
authorized clients in their use of the CAPMC service.

System power calls backed by data from the PMDB, are also subject to the
same hysteresis constraint as the node energy queries.

System power queries can be resource intensive. Depending on system size
and input parameters, those calls may require several minutes to
complete.

### <span class="nh-number">4.4.5. </span>Emergency Power Off

In the event of a catastrophic event where the Cray EX liquid-cooled
hardware is not shutting down automatically, a software driven Emergency
Power Off (EPO) is provided. This causes a hard power off for all
targeted components and their children. The power is dropped at the
Chassis level and there is no clean shutdown for any of the components
below the targeted chassis and related chassis. This mechanism should
not be used for routine power control of the system. The Cray EX CDU,
CMM, and non-Olympus hardware will not be impacted and stay powered on.
EPO is only available at the system, cabinet, and chassis component
levels for Olympus hardware. An EPO on a single chassis will result in
an EPO being sent to every chassis in a cooling group. Each cooling
group may have up to 6 cabinets containing 48 chassis.

# <span class="nh-number">5. </span>Future Shasta Power Management

Over the next few releases, CAPMC functionality is going to be replaced.
Instead, a dedicated microservice for power control, Power Control
Service (PCS), will take over most CAPMC functionality with an idiomatic
REST interface. Remaining features will be split between the Hardware
State Manager (HSM), and the Shasta Monitoring Framework (SMF). PCS will
be the authoritative API for power state and power control while other
services will provide long term trending and inventory management
features. Additionally, PCS will implement the p-state/c-state controls
that are currently missing from CAPMC.

  

## Comments:

<table data-border="0" width="100%">
<colgroup>
<col style="width: 100%" />
</colgroup>
<tbody>
<tr class="odd">
<td><span id="comment-205033925"></span>
<p>This reads like a very technical document for an audience that is already familiar with the challenges and solutions involved.  I think we need a preamble or introduction that sets context for the rest of the paper.  I'm looking for a paragraph that clearly states who the paper is for, what we assume they already know, and what we hope they will learn or be able to do once they have digested the content.</p>
<p>To Andrew's point, we need to talk about generic features separately from the CAPMC implementation of those features.  Describing what CAPMC does in 1.4 today is good, but we also need to reference that Shasta 1.5 will maintain the same functionality, but in different microservices.</p>
<p>I'm guessing that power on and power off are pretty common, but the power capping feature was new to me when I joined HPC.  Depending on who is reading this, it may not be something they are familiar with either.  We should spend more time in the introduction the hardware and software that allow CSM to set upper limits and estimate minimum draw on a per-node bases and extrapolate to per-cabinet and per-system power numbers.</p>
<div class="smallfont" data-align="left" style="color: #666666; width: 98%; margin-bottom: 10px;">
<img src="images/icons/contenttypes/comment_16.png" width="16" height="16" /> Posted by alovelltr at Apr 13, 2021 11:13
</div></td>
</tr>
<tr class="even">
<td style="border-top: 1px dashed #666666"><span id="comment-205033936"></span>
<p>Then I need help. I write short, sweet, technical documentation not white papers.</p>
<p>I thought this was v1.4 focused. The Future section talks briefly about the changes coming with PCS and how to get current CAPMC information down the road.</p>
<div class="smallfont" data-align="left" style="color: #666666; width: 98%; margin-bottom: 10px;">
<img src="images/icons/contenttypes/comment_16.png" width="16" height="16" /> Posted by mjendrysik at Apr 13, 2021 11:21
</div></td>
</tr>
</tbody>
</table>

Document generated by Confluence on Jan 14, 2022 07:17

[Atlassian](http://www.atlassian.com/)
