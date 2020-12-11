.. Copyright 2015,2016 Cray Inc. All Rights Reserved.

Node Capabilities & Power Control
=================================

CAPMC power capping controls implement a simple interface for querying
component capabilities and manipulation of node or sub-node (accelerator)
power constraints. This functionality enables external software to establish
an upper bound, or estimate a minimum bound, on the amount of power a system
or a select subset of the system may consume. The following API calls are
provided as a means for third party software to implement advanced power
management strategies.

Several calls do not include a command line interface. A caller may utilize
the **capmc json** functionality to send and receive customized JSON data
structures conforming to the specification for these calls.

The API calls, **get_power_bias**, **set_power_bias**, **clr_power_bias**,
**set_power_bias_data**, & **compute_power_bias**, are new in SMW-8.0.


get_power_cap_capabilities
--------------------------
.. http:post:: /capmc/get_power_cap_capabilities
   :synposis: Return power capping capabilities

   Get power capping capabilities.

   The **get_power_cap_capabilities** call informs third-party software
   about installed hardware and its associated properties. Information
   returned includes the specific hardware types, NID membership, and
   power capping controls along with their allowable ranges. Information
   may be returned for a targeted set of NIDs or the system as a whole.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_power_cap_capabilities HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 71
      Content-Type: application/json

      {
        "nids": [
          40, 
          41, 
          42, 
          43, 
          48, 
          49
        ]
      }

   :<json array[int] nids: User specified list, or empty array for all NIDs. This
        list must not contain invalid or duplicate NID numbers. If invalid
        NID numbers are specified then an error will be returned.


   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "", 
        "groups": [
          {
            "name": "01:000d:206d:0073:0008:0020:3a34:8100", 
            "desc": "ComputeANC_SNB_115W_8c_32GB_14900_KeplerK20XAccel", 
            "host_limit_max": 185, 
            "host_limit_min": 95, 
            "static": 0, 
            "supply": 425,
            "powerup": 120, 
            "nids": [
              48, 
              49
            ], 
            "controls": [
              {
                "name": "accel",
                "desc": "Accelerator control", 
                "max": 225, 
                "min": 180 
              }, 
              {
                "name": "node",
                "desc": "Node manager control", 
                "max": 410, 
                "min": 275 
              }
            ] 
          }, 
          {
            "name": "01:000d:206d:0104:0010:0040:3200:0000", 
            "desc": "ComputeANC_SNB_260W_16c_64GB_12800_NoAccel", 
            "host_limit_max": 350, 
            "host_limit_min": 200, 
            "static": 0, 
            "supply": 425,
            "powerup": 150,
            "nids": [
              40, 
              41, 
              42, 
              43
            ], 
            "controls": [
              {
                "name": "node",
                "desc": "Node manager control", 
                "max": 350, 
                "min": 200 
              }
            ]
          }
        ]
      }


   :>json int e: Request status code, zero on success

   :>json string err_msg: Human readable error message

   :>json array[object] groups: Object array containing hardware specific
        information and NID membership, each element represents a unique
        hardware type 

   :>json string groups.name: Opaque identifier which Cray system management
        software uses to uniquely identify a node type

   :>json string groups.desc: Text description of the opaque node type
        identifier

   :>json int groups.host_limit_max: Estimated maximum power, specified in
        watts, which host CPU(s) and memory may consume

   :>json int groups.host_limit_min: Estimated minimum power, specified in
        watts, which host CPU(s) and memory require to operate

   :>json int groups.static: Static per node power overhead, specified in
        watts, which is unreported

   :>json int groups.supply: Maximum capacity of each node level power
        supply for the given hardware type, specified in watts

   :>json int groups.powerup: Typical power consumption of each node
        during hardware initialization, specified in watts

   :>json array[int] groups.nids: NID members belonging to the given hardware
        type

   :>json array[object] groups.controls: Array of node level control objects which
        may be assigned or queried, one element per control

   :>json string groups.controls.name: Unique control object identifier

   :>json string groups.controls.desc: Human readable description of the
        control object

   :>json int groups.controls.min: Minimum value which may be assigned to
        the control object, units are dependent upon control type

   :>json int groups.controls.max: Maximum value which may be assigned to
        the control object, units are dependent upon control type

   :status 200: Network API call success
   :status 500: Internal command failure


get_power_cap
-------------
.. http:post:: /capmc/get_power_cap
   :synopsis: Return power capping controls

   Get power capping controls.

   The **get_power_cap** call returns the power capping control(s) and
   currently applied settings for the requested list of NIDs. Control
   values which are returned as zero have special meaning. In such case, a
   zero value indicates the respective control is unconstrained.

   The request must **POST** a properly formatted JSON object to
   the API server. The command takes a single argument which
   identifies a NID list.

   **Example Request**:


   .. sourcecode:: http

      POST /capmc/get_power_cap HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 35
      Content-Type: application/json

      {
        "nids": [
          40, 
          48, 
        ]
      }

   :<json array[int] nids: User specified list, or empty array for all NIDs. This
        list must not contain invalid or duplicate NID numbers. If invalid
        NID numbers are specified then an error will be returned.
        If empty, the default is all NIDs. The specified NIDs must be
        in the **ready** state per the **node_status** command.


   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "", 
        "nids": [
          {
            "nid": 40,
            "controls": [
              {
                "name": "node", 
                "val": 350
              },
              {
                "name": "node-biased", 
                "val": 355
              }, 
              {
                "name": "bias-factor", 
                "val": 1.015471
              }
            ] 
          }, 
          {
            "nid": 48,
            "controls": [
              {
                "name": "node", 
                "val": 0
              }, 
              {
                "name": "accel", 
                "val": 0
              }
            ] 
          }
        ]
      }


   **Example Response**: (Partial Success or Failure)

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 52, 
        "err_msg": "Invalid exchange", 
        "nids": [
          {
            "nid": 40,
            "e": 52,
            "err_msg": "Invalid exchange"
          }, 
          {
            "nid": 48,
            "controls": [
              {
                "name": "node", 
                "val": 0
              }, 
              {
                "name": "accel", 
                "val": 0
              }
            ] 
          }
        ]
      }


   :>json int e: Overall request status code, zero on total success, non-zero
        if one or more node specific operations fail

   :>json string err_msg: Human readable error message

   :>json array[object] nids: Object array containing NID specific result data,
        each element represents a single NID 

   :>json int nids.nid: NID number owning the returned control objects

   :>json int nids.e: Optional, error status, non-zero indicates
        operation failed on this node

   :>json string nids.err_msg: Optional, human readable error message applicable
        to this node

   :>json array[object] nids.controls: Optional, array of node level control
        and status objects which have been queried, one element per control

   :>json string nids.controls.name: Unique control or status object
        identifier

   :>json int nids.controls.val: Control object setting, or zero to
        indicate control is unconstrained, units are dependent upon
        control type

   :status 200: Network API call success
   :status 500: Internal command failure


set_power_cap
-------------
.. http:post:: /capmc/set_power_cap
   :synopsis: Set power capping parameters

   Set power capping parameters.

   The **set_power_cap** call is used to establish an upper bound with
   respect to power consumption on a per-node, and if applicable, a
   sub-node basis.  Established power cap parameters will revert to the
   default configuration on the next system boot.

   If setting multiple different power caps is desired, then it is
   recommended that those be set programmatically via the **json** command
   with an input data structure conforming to the HTTP Interface request
   format for this command.  The **json** command allows third party
   software to pass its own JSON formatted requests in a single
   transaction to the HTTP API service.

   Service nodes may not be power capped. If service node NIDs are
   specified then the request will fail with an invalid parameters error.
   When applying a power cap, unspecified controls are reset to their
   default value.


   Specify the NIDs to apply the specified power caps. The syntax allows a
   comma-separated list of nids (eg, "1,4,5"), a range of nids (eg, "7-10"), or
   both (eg, "1,4,5,7-10").

   Nodes with high powered accelerators and high TDP processors will be
   automatically power capped at the "supply" limit returned per the
   **get_power_cap_capabilities** command. If a node level power cap is
   specified that is within the node control range but exceeds the supply
   limit, the actual power cap assigned will be clamped at the supply limit.

   Specify the desired accelerator component power cap. The value given must
   be within the range returned in the capabilities output. A value of zero
   may be supplied to to explicitly clear an accelerator power cap.

   The accelerator power cap value represents a subset of the total node level
   power cap. If a node level power cap of 400 watts is applied and an
   accelerator power cap of 180 watts is applied, then the total node power
   consumption is limited to 400 watts. If the accelerator is actively
   consuming its entire 180 watt power allocation, then the host processor,
   memory subsystem, and support logic for that node may consume a maximum of
   220 watts.


   In the common case, the response payload is short and consists only of an
   integer status code and an optional message. However there may be
   instances, likely due to hardware errors, where a small number of nodes
   encounter a problem and are unable to comply with the command. If an error
   does occur, extra information pertaining to the specific component where
   the error occurred is included in the response payload.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes an array of objects which identify target
   component NIDs, control names, and their associated set point values. 

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/set_power_cap HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 472
      Content-Type: application/json

      {
        "nids": [
          {
            "nid": 20,
            "controls": [
              {
                "name": "node", 
                "val": 400
              }
            ]
          }, 
          {
            "nid": 21,
            "controls": [
              {
                "name": "node", 
                "val": 400
              }
            ]
          }, 
          {
            "nid": 60,
            "controls": [
              {
                "name": "node", 
                "val": 410
              }, 
              {
                "name": "accel", 
                "val": 220
              }
            ]
          }
        ]
      }

   :<json array[object] nids: Object array containing NID specific input
        data, each element represents a single NID.
	The specified NIDs must be in the **ready** state per the
	**node_status** command.

   :<json int nids.nid: NID number owning the specified input control objects

   :<json object[] nids.controls: Array of node level control objects to be
        adjusted, one element per control

   :<json string nids.controls.name: Unique control object identifier
   :<json int nids.controls.val: Control object setting, or zero to
        indicate control is unconstrained, units are dependent upon control
        type.


   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
          "e":0,
          "err_msg":""
      }

   **Example Response**: (Partial Success or Failure)

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 52, 
        "err_msg": "Invalid exchange", 
        "nids": [
          {
            "nid": 60,
            "e": 52, 
            "err_msg": "Invalid exchange" 
          }
        ]
      }


   :>json int e: Request status code, zero on success.

   :>json string err_msg: Human readable error message
   
   :>json object[] nids: Object array containing NID specific error data, NIDs
        which experienced success are omitted
   
   :>json int nids[].nid: NID number owning the returned error data
   
   :>json int nids[].e: Error status, non-zero indicates operation failed on
        this node
   
   :>json string nids[].err_msg: Human readable error string applicable to this
        node

   :status 200: Network API call success
   :status 500: Internal command failure



get_power_bias
--------------
.. http:post:: /capmc/get_power_bias
   :synopsis: Return power bias

   Get poer bias.

   The **get_power_bias** API call informs third-party software what, if
   any, per node multiplication factor will be considered by low level HSS
   software when applying a node level power cap. When low level HSS
   assigns a node level power cap, it assigns the product of the caller
   specified value and the per node bias factor as the actual power cap.
   By default, each node is assigned a power bias factor of 1.0. This
   results in the actual power cap being equal to the caller specified
   value unless a power bias factor has been explicitly configured.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_power_bias HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 35
      Content-Type: application/json

      {
        "nids": [
          40, 
          41
        ]
      }

   :<json array[int] nids: User specified list, or empty array for all NIDs. This
        list must not contain invalid or duplicate NID numbers. If invalid
        NID numbers are specified then an error will be returned.

   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success", 
        "nids": [
          {
            "nid": 41, 
            "power-bias": 1.0
          }, 
          {
            "nid": 40, 
            "power-bias": 1.0
          }
        ]
      }
      
   :>json int e: Overall request status code, zero on total success, non-zero
        if one or more node specific operations fail
   :>json string err_msg: Human readable error message
   :>json array[object] nids: Object array containing NID specific result data,
        each element represents a single NID 
   :>json int nids.nid: NID number owning the returned power-bias value
   :>json float nids.power-bias: Power bias setting, or 1.0 to indicate the
        default
        
   :status 200: Network API call success
   :status 500: Internal command failure


set_power_bias
--------------
.. http:post:: /capmc/set_power_bias
   :synopsis: Set node power bias

   Set power bias.

   A caller may establish a per node power capping bias factor via the
   **set_power_bias** API call. This may be used as a fine grained tuning
   knob intended to equalize node to node performance variation, through
   the dithering of individual node level power caps, while operating the
   system under a global power cap. A caller may derive the power cap bias
   factors by any means, or use the built in **set_power_bias_data** and
   **compute_power_bias** API calls.

   Newly established power cap bias factors do not take effect until the
   respective node level power cap has been reapplied.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a list of objects where each element contains two
   values which identify a specific NID and an associated power bias value.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/set_power_bias HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 139
      Content-Type: application/json

      {
        "nids": [
          {
            "nid": 41, 
            "power-bias": 0.984529
          }, 
          {
            "nid": 42, 
            "power-bias": 1.015471
          }
        ]
      }

   :<json array[object] nids: Object array containing NID specific input data,
        each element represents a single NID 
   :<json int nids.nid: NID number owning the assigned power-bias value
   :<json float nids.power-bias: Power bias setting, or 1.0 to reset to
        default

   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success"
      }

   :>json int e: Request status code, zero on success
   :>json string err_msg: Human readable error message

   :status 200: Network API call success
   :status 500: Internal command failure


clr_power_bias
--------------
.. http:post:: /capmc/clr_power_bias
   :synopsis: Clear per node power capping bias

   Clear node power capping bias factors

   A caller may clear per node power capping bias factors via the
   **clr_power_bias** API call. This call differs from calling
   **set_power_bias** with a specified power bias of 1.0 in that this call
   results in all internal records relating to the assigned power bias
   being deleted.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/clr_power_bias HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length:35 
      Content-Type: application/json

      {
        "nids": [
          41,
          42
        ]
      }

   :<json array[int] nids: User specified list, or empty array for all NIDs. 

   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success"
      }

   :>json int e: Request status code, zero on success
   :>json string err_msg: Human readable error message

   :status 200: Network API call success
   :status 500: Internal command failure


set_power_bias_data
-------------------
.. http:post:: /capmc/set_power_bias_data
   :synposis: Set power capping bias data

   Set power bias data.

   Average power per NID over an application run may be stored for later
   processing using the **set_power_bias_data** API call. This information
   is primarily used within the **compute_power_bias** API call. It is
   intended that this call be used as part of a higher level system
   characterization process.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes an application name and list of objects where
   each element contains two values which identify a specific NID and an
   associated average power value.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/set_power_bias_data HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 141
      Content-Type: application/json

      {
        "app": "stress", 
        "nids": [
          {
            "avgpwr": 350, 
            "nid": 41
          }, 
          {
            "avgpwr": 361, 
            "nid": 42
          }
        ]
      }

   :<json array[object] nids: Object array containing NID specific input data,
        each element represents a single NID 
   :<json int nids.nid: NID number owning the assigned power-bias value
   :<json int nids.avgpwr: Average power consumption over the application
        run, specified in watts


   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success"
      }

   :>json int e: Request status code, zero on success
   :>json string err_msg: Human readable error message

   :status 200: Network API call success
   :status 500: Internal command failure


compute_power_bias
------------------
.. http:post:: /capmc/compute_power_bias
   :synopsis: Compute power bias

   Compute power bias.

   The **compute_power_bias** API call is used to calculate a per node
   power cap multiplication factor for each NID as it relates to a larger
   set of NIDs for a given application. The computed values returned by
   this API call are not automatically saved. If so desired, the values
   must be explicitly saved using the **set_power_bias** API call. Prior
   to using this API call, the caller must have previously stored average
   power data for the specified application and target NID list using the
   **set_power_bias_data** API call.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes two arguments which identify a NID list and
   application name.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/compute_power_bias HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 55
      Content-Type: application/json

      {
        "app": "stress", 
        "nids": [
          41, 
          42
        ]
      }

   :<json string app: Application name
   :<json array[int] nids: User specified list, or empty array for all NIDs in which
        an average application power record for the specified application
        exists


   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success", 
        "nids": [
          {
            "nid": 41, 
            "power-bias": 0.984529
          }, 
          {
            "nid": 42, 
            "power-bias": 1.015471
          }
        ]
      }


   :>json int e: Overall request status code, zero on total success, non-zero
        if one or more node specific operations fail
   :>json string err_msg: Human readable error message
   :>json array[object] nids: Object array containing NID specific result data,
        each element represents a single NID 
   :>json int nids.nid: NID number owning the returned power-bias value
   :>json float nids.power-bias: Power bias setting

   :status 200: Network API call success
   :status 500: Internal command failure

