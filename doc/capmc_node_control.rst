.. Copyright 2015-2018 Cray Inc. All Rights Reserved.

Node Status & On Off Control
============================

CAPMC node power controls implement a simple interface for powering on or off
compute nodes, querying node state information, and querying site-specific
service usage rules. These controls enable external software to more
intelligently manage system wide power consumption or configuration
parameters. The simplest power management strategy may be to simply turn off
compute nodes which may be idle for a significant time interval and turn them
back on when demand increases. The following API calls are provided as a means
for third party software to implement power management strategies as simple or
complex as the site-level requirements demand.


node_rules
----------

.. http:post:: /capmc/get_node_rules
   :synopsis: Returns node rules 

   Get node rules.

   The **node_rules** API informs third party software about hardware (and
   perhaps site-specific) rules or timing constraints that allow for
   efficient and effective management of idle node resources. The data
   returned informs the caller of how long **node_on** and **node_off**
   operations are expected to take, the minimum amount of time nodes
   should be left off to save energy, and limits on the number of nodes
   that should be turned on or off at once. Default rules are supplied
   where appropriate.

   Other values such as the maximum node counts for **node_on** or
   **node_off** and the maximum amount of time a node should remain off
   after a power down are left unset. The values are not strictly enforced
   by Cray system management software. They are meant to provide
   guidelines for authorized callers in their use of the CAPMC service.

   The request must **POST** an empty JSON object to the API server.  This
   command takes no arguments.

   **Example request**:

   .. sourcecode:: http

      POST /capmc/get_node_rules HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 2
      Content-Type: application/json

      {}


   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "", 
        "latency_node_off": 60, 
        "latency_node_on": 600, 
        "latency_node_reinit": 760,
        "max_off_req_count": -1, 
        "max_off_time": -1, 
        "max_on_req_count": -1, 
        "max_reinit_req_count": -1,
        "min_off_time": 900
      }
      
   :>json int e: Request status code, zero on success

   :>json string err_msg: Human readable error message

   :>json int latency_node_off: Approximate time, in seconds, in which a node
        is expected to perform a clean shutdown and power off

   :>json int latency_node_on: Approximate time, in seconds, in which a node
        is expected to perform a warm bounce and boot into a ready state

   :>json int latency_node_reinit: Approximate time, in seconds, in which a node
        is expected to perform a clean shutdown, power off, warm bounce, and
        boot into a ready state

   :>json int max_off_req_count: Maximum number of nodes which may be powered
        off at once. -1 indicates no limit

   :>json int max_off_time: Maximum time, in seconds, in which a node may be 
        in the powered off state. -1 indicates no limit

   :>json int max_on_req_count: Maximum number of nodes which may be powered
        on and booted at once. -1 indicates no limit

   :>json int max_reinit_req_count: Maximum number of nodes which may be
        shutdown, powered off, powered on, and rebooted at once.  -1 indicates
        no limit

   :>json int min_off_time: Minimum time, in seconds, in which a node must
        remain in the powered off state after a shutdown and power off
        operation. -1 indicates no limit

   :status 200: Success


.. note::

   Default rules are established automatically at system installation
   time.  The administrator may choose to customize the rule set.
   Customization is performed by editing the respective parameter in a
   configuration file (**/opt/cray/hss/default/etc/xtremoted/rules.ini**)
   stored on the System Management Workstation.


node_status
-----------

.. http:post:: /capmc/get_node_status
   :synopsis: Returns node status 

   Get node status

   A node's component state may be returned via the **node_status**
   function for the full system or a subset as specified by a nid list or
   component filter. The status API call is intended, but not limited, to
   be used in conjunction with asynchronous operations which may modify
   node component state, such as **node_on** or **node_off**.

   Third party utilities would issue an asynchronous operation, such as
   **node_on**, and if the operation was successfully enqueued poll for
   changes in component state after waiting for the expected boot time
   latency.  If the targeted component state has switched from "off" to
   "ready" then the caller knows the operation was successful.

   States reported through this API call mirror those defined in Cray HSS.

   **Node States:**

       * **disabled** - Component is physically installed,
         but ignored by Cray system management software 
       * **halt** - Operating system has shut down, hardware has not yet powered off
       * **on** - Power is on and BIOS has initialized all hardware, node is waiting to be booted
       * **off** - Power is off
       * **ready** - Operating system is fully booted
       * **standby** - Operating system is in process of booting


      **Example Request**:

      .. sourcecode:: http

         POST /capmc/get_node_status HTTP/1.1
         Host: smw.example.com
         Accept: application/json
         Content-Length: 96
         Content-Type: application/json

         {
           "filter": "show_on|show_ready", 
           "nids": [
             1, 
             40, 
             41, 
             42, 
             43
           ]
         }

   :<json string filter: Pipe concatenated (|) list of filter strings
   :<json array[int] nids: User specified list, or empty array for all
    NIDs. This list must not contain invalid or duplicate NID numbers.
    If invalid NID numbers are specified then an error will be returned.

    **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "", 
        "on": [
          40, 
          42, 
          43, 
          41
        ], 
        "ready": [
          1
        ]
      }

   :>json int e: Request status code, zero on success, nonzero on invalid input.

   :>json string err_msg: Human readable error message

   :>json array[int] disabled: Optional member, list of disabled NIDs

   :>json array[int] halt: Optional, list of halted NIDs

   :>json array[int] on: Optional, list of powered on NIDs

   :>json array[int] off: Optional, list of powered off NIDs

   :>json array[int] ready: Optional, list of booted NIDs

   :>json array[int] standby: Optional, list of booting NIDs

   :status 200: Success
   :status 400: Bad Request
   :status 500: Internal Server Error

.. note::

     The **get_node_status** API call does not report **empty** components.

     Specify filters for status query. Filters may be pipe-separated (|)
     and surrounded with double quotes, e.g.  "opt|opt|opt". Valid filters
     are: **show_all**, **show_off**, **show_on**, **show_halt**,
     **show_standby**, **show_ready**, **show_diag**, and
     **show_disabled**. If omitted, the default is **show_all**.

     The request must **POST** a properly formatted JSON object to the API
     server. The command takes two optional arguments which identify a NID
     list and component status filter.

node_on
-------

.. http:post:: /capmc/node_on
   :synopsis: Power on nodes

   Power on and boot a list of nodes.

   Power on and warm boot a selected list of compute node NIDs using the
   default boot image. This has no effect on the status of the high speed
   network (HSN).  However, this command requires that the HSN ASIC
   attached to each node in the target list has previously been powered on
   and routed.

   The **node_on** API call is **asynchronous**. It returns immediately
   containing a status result indicating an error with invalid input
   parameters, or success indicating the operation has been enqueued into
   an asynchronous command processing queue. The caller must determine
   overall command status by polling for **node_status** after the
   expected power on and warm boot period has elapsed.

   The request must **POST** an array of selected NIDs along with an optional
   human readable reason string.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/node_on HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 85
      Content-Type: application/json

      {
        "reason": "need more capacity", 
        "nids": [
          40, 
          41, 
          42, 
          43
        ]
      }

   :<json string reason: Arbitrary, free-form text
   :<json array[int] nids: User specified list of compute node NIDs to warm bounce
        and boot. An empty array is invalid. If invalid NID numbers are
        specified then an error will be returned.

   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": ""
      }

   :>json int e: Request status code, zero on success, nonzero on invalid input.
   :>json string err_msg: Human readable error message

   :status 200: Success
   :status 400: Bad Request

node_off
--------

.. http:post:: /capmc/node_off
   :synopsis: Power off nodes

   Shudown and power off nodes.

   Shutdown and power off a selected list of compute node NIDs. This has
   no effect on the status of the high speed network (HSN). The HSN ASIC
   attached to each node will remain powered on and routed. After the
   **node_off** operation has completed, the selected nodes will be in a
   state suitable for warm booting back into the system at a later date.

   The **node_off** API call is **asynchronous**. It returns immediately
   containing a status result indicating an error with invalid input
   parameters, or success indicating the operation has been enqueued into
   an asynchronous command processing queue. The caller must determine
   overall command status by polling for **node_status** after the
   expected shutdown and power off period has elapsed.

   Specify an arbitrary text message which is given as the reason for
   performing the node_off operation. This argument is optional and is used in
   the same manner as with the **node_on** command.

   The request must **POST** an array of selected NIDs along with an optional
   human readable reason string.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/node_off HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 96
      Content-Type: application/json

      {
        "reason": "powersave, need less capacity", 
        "nids": [
          40, 
          41, 
          42, 
          43
        ]
      }

   :<json string reason: Arbitrary, free-form text
   :<json array[int] nids: User specified list of compute node NIDs to shutdown and
        power off. An empty array is invalid. If invalid NID numbers are
        specified then an error will be returned.

   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": ""
      }

   :>json int e: Request status code, zero on success, nonzero on invalid input.
   :>json string err_msg: Human readable error message

   :status 200: Success
   :status 400: Bad Request


node_reinit
-----------

.. http:post:: /capmc/node_reinit
   :synopsis: Power off and reboot nodes

   Shutdown and reboot nodes.
   
   Warmbounce and boot a selected list of compute node NIDs using the
   default boot image. This has no effect on the status of the high speed
   network (HSN).  However, this command requires that the HSN ASIC
   attached to each node in the target list has previously been powered on
   and routed.

   The **node_reinit** API call is **asynchronous**. It returns
   immediately containing a status result indicating an error with invalid
   input parameters, or success indicating the operation has been enqueued
   into an asynchronous command processing queue.  The caller must
   determine overall command status by polling for **node_status** after
   the expected warmbounce and boot period
   has elapsed.

   Specify an arbitrary text message which is given as the reason for
   performing the node_off operation. This argument is optional and is
   used in the same manner as with the **node_on** command.

   The request must **POST** an array of selected NIDs along with an
   optional human readable reason string.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/node_reinit HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 102
      Content-Type: application/json

      {
        "reason": "apply staged node configurations", 
        "nids": [
          40, 
          41, 
          42, 
          43
        ]
      }

   :<json string reason: Arbitrary, free-form text
   :<json array[int] nids: User specified list of compute node NIDs to shutdown and
        power off. An empty array is invalid. If invalid NID numbers are
        specified then an error will be returned.

   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": ""
      }

   :>json int e: Request status code, zero on success, nonzero on invalid input.
   :>json string err_msg: Human readable error message

   :status 200: Success
   :status 400: Bad Request
