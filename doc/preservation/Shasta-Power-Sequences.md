1.  [CASMHMS](index.html)
2.  [CASMHMS Home](CASMHMS-Home_119901124.html)
3.  [Design Documents](Design-Documents_127906417.html)
4.  [CAPMC](CAPMC_227479290.html)

# <span id="title-text"> CASMHMS : Shasta Power Sequences </span>

Created by <span class="author"> Matt Kelly</span>, last modified on May
07, 2019

This page outlines the overall sequencing of power operations starting
with a cold, powered off system.  It is not meant to be a
customer-facing document, but rather the basis for one.

-   <span class="TOCOutline">1</span>
    [Scope](#ShastaPowerSequences-Scope)
-   <span class="TOCOutline">2</span> [System
    Description](#ShastaPowerSequences-SystemDescription)
-   <span class="TOCOutline">3</span>
    [Breakers](#ShastaPowerSequences-Breakers)
    -   <span class="TOCOutline">3.1</span> [What Is Powered
        Up](#ShastaPowerSequences-WhatIsPoweredUp)
    -   <span class="TOCOutline">3.2</span> [Status Upon
        Completion](#ShastaPowerSequences-StatusUponCompletion)
-   <span class="TOCOutline">4</span> [NCN and SMS
    Nodes](#ShastaPowerSequences-NCNandSMSNodes)
    -   <span class="TOCOutline">4.1</span> [What Is Powered
        Up](#ShastaPowerSequences-WhatIsPoweredUp.1)
    -   <span class="TOCOutline">4.2</span> [Status Upon
        Completion](#ShastaPowerSequences-StatusUponCompletion.1)
-   <span class="TOCOutline">5</span> [River Compute Nodes (For
    Discovery)](#ShastaPowerSequences-RiverComputeNodes(ForDiscovery))
    -   <span class="TOCOutline">5.1</span> [What Is Powered
        Up](#ShastaPowerSequences-WhatIsPoweredUp.2)
    -   <span class="TOCOutline">5.2</span> [Status Upon
        Completion](#ShastaPowerSequences-StatusUponCompletion.2)
-   <span class="TOCOutline">6</span>
    [Storage](#ShastaPowerSequences-Storage)
-   <span class="TOCOutline">7</span> [Mountain Chassis
    Rectifiers](#ShastaPowerSequences-MountainChassisRectifiers)
    -   <span class="TOCOutline">7.1</span> [What Is Powered
        Up](#ShastaPowerSequences-WhatIsPoweredUp.3)
    -   <span class="TOCOutline">7.2</span> [Status Upon
        Completion](#ShastaPowerSequences-StatusUponCompletion.3)
-   <span class="TOCOutline">8</span> [Mountain Switch
    Blades](#ShastaPowerSequences-MountainSwitchBlades)
    -   <span class="TOCOutline">8.1</span> [What Is Powered
        Up](#ShastaPowerSequences-WhatIsPoweredUp.4)
    -   <span class="TOCOutline">8.2</span> [Status Upon
        Completion](#ShastaPowerSequences-StatusUponCompletion.4)
-   <span class="TOCOutline">9</span> [Mountain Compute
    Blades](#ShastaPowerSequences-MountainComputeBlades)
    -   <span class="TOCOutline">9.1</span> [What Is Powered
        Up](#ShastaPowerSequences-WhatIsPoweredUp.5)
    -   <span class="TOCOutline">9.2</span> [Status Upon
        Completion](#ShastaPowerSequences-StatusUponCompletion.5)
-   <span class="TOCOutline">10</span> [River Compute
    Nodes](#ShastaPowerSequences-RiverComputeNodes)
    -   <span class="TOCOutline">10.1</span> [What Is Powered
        Up](#ShastaPowerSequences-WhatIsPoweredUp.6)
    -   <span class="TOCOutline">10.2</span> [Status Upon
        Completion](#ShastaPowerSequences-StatusUponCompletion.6)
-   <span class="TOCOutline">11</span> [Mountain Compute
    Nodes](#ShastaPowerSequences-MountainComputeNodes)
    -   <span class="TOCOutline">11.1</span> [What Is Powered
        UP](#ShastaPowerSequences-WhatIsPoweredUP)
    -   <span class="TOCOutline">11.2</span> [Status Upon
        Completion](#ShastaPowerSequences-StatusUponCompletion.7)

## Scope

This document will cover power-up sequences from "ground zero".   This
means that a system is completely powered off – no power even to any
racks or cabinets, no CDUs running, nothing at all.

One assumption that is being made: all hardware that is configurable, be
it firmware flashing, BIOS initial settings, front-panel display
settings, etc. has already been done.  This guide does not include any
of those steps.

Also outside the scope of this document are any details of SW
installation or configuration.

## System Description

The power-up sequences will include steps for a fully-equipped Shasta
system, including:

-   System Management Services NCNs
-   User-access NCNs
-   River compute and SSN nodes and racks
-   River CDUs (standalone and in-rack)
-   Mountain compute cabinets
-   Mountain CDUs
-   Rosetta HSN
-   System Management (HW mgmt and node mgmt) ethernet networks
-   Storage

## Breakers

The first step is to power up all River and Mountain racks/cabinets/CDUs
via the wall breakers.  This is a manual step.

### What Is Powered Up

Once this is done the following HW will be powered up:

-   River  
    -   River rack PDU and all sockets
    -   Management Network ethernet switches
    -   Rosetta TOR switches
    -   River blade BMCs, includes compute nodes and all NCNs.
-   Storage
    -   Storage controller blade BMCs
    -   Storage disks/SSDs
-   Mountain
    -   Cabinet CECs and front panel displays
    -   Chassis-level CMMs

### Status Upon Completion

-   The BMCs of all River nodes, and the CMMs of Mountain cabinets are
    available for command/control operations
-   SMS nodes can be powered on and installed/configured. 
-   River nodes can be turned on
-   Mountain CDU cooling can be enabled via component power-on
-   Mountain chassis-level rectifiers can be turned on
-   Node and switch blades can be powered on

## NCN and SMS Nodes

The next step is to power on NCN nodes.  These include UAN nodes but
most importantly, SMS nodes, since all  L2/L3 management software
requires SMS.  Typically these are powered up manually and are either
discovered by MAAS or via REDS.

### What Is Powered Up

-   SMS nodes (can be staged; first a single one, then more as needed
    when SW install progresses)
-   UAN nodes (can be done now or later)

### Status Upon Completion

-   SMS and UAN nodes (if powered on) are available to have some level
    of SW installed on them.  Most likely a full SMS install is
    performed.

## River Compute Nodes (For Discovery)

At this point the River compute nodes are turned on.  This allows them
to be discovered via REDS.  Currently this is a manual step. 

### What Is Powered Up

-   River compute blades.

### Status Upon Completion

-   River blade BMC MAC addresses are known (for DHCP and further HW
    discovery)
-   River nodes are powered down (done by the REDS process)
-   River blade BMC credentials have been set and are known
-   All River blade BMC information is in the Hardware State Manager

At this point the HSM can be instructed to do a HW discovery operation,
during which the BMCs are queried for various information about the
River nodes (MAC addresses, etc.).  This information is stored in the
HSM database.

## Storage

Any Kilamanjaro storage racks associated with the system contain
controller nodes which are at this point controllable via their BMCs but
not yet turned on. 

ZZZZ

## Mountain Chassis Rectifiers

Mountain chassis rectifiers must now be turned on in order to supply
power for compute and switch blades.

### What Is Powered Up

-   Mountain Chassis Rectifiers

### Status Upon Completion

-   Chassis Rectifiers are powered on
-   Node and Switch blades are NOT powered on

## Mountain Switch Blades

Once the chassis rectifiers are enabled, Mountain switch blades can be
turned on. 

### What Is Powered Up

-   Mountain switch blades

### Status Upon Completion

-   Mountain switch blades are powered up
-   Switch blades' HSN hardware is powered up
-   Rosetta TORs are turned on, as are now the Mountain switch blades;
    the HSN could theoretically be initialized.

The next steps involve turning on compute nodes.  These nodes will then
boot OS.  Thus, the Boot Script Service needs to be ready to dole out
boot images to the nodes about to be booted; boot orchestration overall
will take over at this point.   Note also that the HSN fabric must be
initialized in order for soon-to-be-booted compute nodes to use the
HSN.  This is done via the HSN fabric and network management SW, which
is outside the scope of this document.

## Mountain Compute Blades

Mountain compute blades can now be powered on.   Note that this does not
turn on the compute nodes, just the blades.

### What Is Powered Up

-   Mountain compute blades (NOT the nodes on the blades).  This is
    analagous to a River blade's BMC powering up but not the node(s)
    controlled by the BMC.

### Status Upon Completion

-   Mountain compute blades are powered on and ready to receive Redfish
    commands to turn on their nodes.

## River Compute Nodes

River compute nodes can now be turned on and thus booted.  Note that
these nodes were turned on once before for the purpose of initial
discovery, and then turned back off.

### What Is Powered Up

-   River compute nodes.   Note that there may be several groups of
    River compute nodes which will be powered on one (or possibly more)
    group at a time.

### Status Upon Completion

-   Powered-up River compute nodes are running their OS.

## Mountain Compute Nodes

Mountain compute nodes are now booted.   Note that this step can be
combined with the previous step, in which case the overall "step" is the
power-on and boot of all compute nodes, be they River or Mountain.

### What Is Powered UP

-   Mountain compute nodes

### Status Upon Completion

-   Powered-up Mountain compute nodes are running their OS

  

  

  

  

  

Document generated by Confluence on Jan 14, 2022 07:17

[Atlassian](http://www.atlassian.com/)
