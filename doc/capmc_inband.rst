.. Copyright 2015,2016 Cray Inc. All Rights Reserved.

Node Frequency & Sleep State Control
====================================

CAPMC node frequency and sleep state controls implement an interface such that
external agents may modify, on the fly, CPU operating frequencies and sleep
state limits. For example, an external agent may reallocate power amongst
nodes by adjusting operating frequencies. A system administrator may disable
package sleep states on idle nodes in an effort to keep system power
consumption inside the power band negotiated with the utility provider. Or,
batch schedulers may optimize power efficiency by running phases of an
application at differing frequencies. The controls are provided as a means for
third-party integrators to implement advanced power management policies.

Runtime P-State and C-State controls are new in SMW-8.0 / CLE-6.0.


get_freq_capabilities
---------------------
.. http:post:: /capmc/cnctl
   :synopsis: Returns supported processor operating frequences

   Get supported processor operating frequencies.

   The **get_freq_capabilities** call informs third-party software about
   supported processor operating frequencies. Valid operating frequencies,
   specified in Hz, are returned in list form.

   The request must **POST** a JSON object to the API server containing a
   remote API call function name, a single element list containing the
   appropriate attribute name, and the protocol version.

   **Example request**:

   .. sourcecode:: http

      POST /capmc/cnctl HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 234
      Content-Type: application/json

      {
        "nids": [
          50
        ],
        "data": {
          "PWR_Function": "PWR_ObjAttrValueCapabilities", 
          "PWR_Attrs": [
            {
              "PWR_AttrName": "PWR_ATTR_FREQ"
            }
          ], 
          "PWR_MajorVersion": 0, 
          "PWR_MinorVersion": 1
        } 
      }

   :<json array[int] nids: User specified NID list, must **not** be empty

   :<json object data: Remote API call input object

   :<json string data.PWR_Function: Remote API call function name, the caller
        must specify "PWR_ObjAttrValueCapabilities"

   :<json array[object] data.PWR_Attrs: Array of attribute name value pairs

   :<json string data.PWR_Attrs.PWR_AttrName: Attribute name, the caller
        must specify "PWR_ATTR_FREQ"

   :<json int data.PWR_MajorVersion: Remote API call major version number,
        the caller must specify "0"

   :<json int data.PWR_MinorVersion: Remote API call minor version number,
        the caller must specify "1"


   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success", 
        "nids": [
          {
            "nid": 50,
            "data": {
              "PWR_Attrs": [
                {
                  "PWR_AttrName": "PWR_ATTR_FREQ", 
                  "PWR_AttrValueCapabilities": [
                    2601000, 
                    2600000, 
                    2500000, 
                    2400000, 
                    2300000, 
                    2200000, 
                    2100000, 
                    2000000, 
                    1900000, 
                    1800000, 
                    1700000, 
                    1600000, 
                    1500000, 
                    1400000, 
                    1300000, 
                    1200000
                  ], 
                  "PWR_ReturnCode": 0
                }
              ], 
              "PWR_ReturnCode": 0,
              "PWR_ErrorMessages": null, 
              "PWR_Messages": null, 
              "PWR_MajorVersion": 0, 
              "PWR_MinorVersion": 1
            }
          }
        ]
      }
 

   :>json int e: Overall request status code, non-zero on partial success or
        failure

   :>json string err_msg: Human readable error message

   :>json array[object] nids: Object array containing NID specific result data, each
        element represents a single NID

   :>json int nids.nid: NID number owning the returned attribute objects

   :>json object nids.data: Remote API call result object

   :>json array[object] nids.data.PWR_Attrs: Array of attribute name value pairs

   :>json string nids.data.PWR_Attrs.PWR_AttrName: Attribute name
        copied from original request

   :>json array[int] nids.data.PWR_Attrs.PWR_AttrValueCapabilities: List of
        supported processor operating frequencies, in Hz

   :>json int nids.data.PWR_Attrs.PWR_ReturnCode: Attribute probe specific
        result code, non-zero on failure

   :>json int nids.data.PWR_ReturnCode: Per-node attribute probe result code,
        non-zero on failure

   :>json string nids.data.PWR_ErrorMessages: Per-node attribute error
        message, or null

   :>json string nids.data.PWR_Messages: Per-node attribute info message, or
        null

   :>json int nids.data.PWR_MajorVersion: Remote API call major version number

   :>json int nids.data.PWR_MinorVersion: Remote API call minor version number

   :status 200: Network API call success
   :status 400: Bad Request
   :status 500: Internal command failure
   :status 504: Gateway Timeout


get_freq_limits
---------------
.. http:post:: /capmc/cnctl
   :synopsis: Return the minium and maxium per-node operating frequences

   Get minimum and maximum allowable per-node operating frequencies.

   The **get_freq_limits** call returns the minimum and maximum allowable
   operating frequencies on a per-node basis. The processor frequency
   operating window may be constrained from defaults using the
   **set_freq_limits** call.

   The request must **POST** a JSON object to the API server containing a
   remote API call function name, a three element list containing the
   appropriate attribute names, and the protocol version.

   **Example request**:

   .. sourcecode:: http

      POST /capmc/cnctl HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 234
      Content-Type: application/json

      {
        "nids": [
          50
        ],
        "data": {
          "PWR_Function": "PWR_ObjAttrGetValues", 
          "PWR_Attrs": [
            {
              "PWR_AttrName": "PWR_ATTR_FREQ"
            }, 
            {
              "PWR_AttrName": "PWR_ATTR_FREQ_LIMIT_MIN"
            }, 
            {
              "PWR_AttrName": "PWR_ATTR_FREQ_LIMIT_MAX"
            }
          ], 
          "PWR_MajorVersion": 0, 
          "PWR_MinorVersion": 1
        }
      }

   :<json array[int] nids: User specified NID list, must **not** be empty

   :<json object data: Remote API call input object

   :<json string data.PWR_Function: Remote API call function name, the caller
        must specify "PWR_ObjAttrGetValues"

   :<json array[object] data.PWR_Attrs: Array of attribute name value pairs

   :<json string data.PWR_Attrs.PWR_AttrName: Attribute name, the caller must
        specify array elements for attributes named "PWR_ATTR_FREQ",
        "PWR_ATTR_FREQ_LIMIT_MIN", and "PWR_ATTR_FREQ_LIMIT_MAX"

   :<json int data.PWR_MajorVersion: Remote API call major version number,
        the caller must specify "0"

   :<json int data.PWR_MinorVersion: Remote API call minor version number,
        the caller must specify "1"


   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success", 
        "nids": [
          {
            "nid": 50,
            "data": {
              "PWR_Attrs": [
                {
                  "PWR_AttrName": "PWR_ATTR_FREQ", 
                  "PWR_AttrValue": 2601000, 
                  "PWR_ReturnCode": 0, 
                  "PWR_TimeNanoseconds": 80864412, 
                  "PWR_TimeSeconds": 1433536488
                }, 
                {
                  "PWR_AttrName": "PWR_ATTR_FREQ_LIMIT_MAX", 
                  "PWR_AttrValue": 2601000, 
                  "PWR_ReturnCode": 0, 
                  "PWR_TimeNanoseconds": 80876654, 
                  "PWR_TimeSeconds": 1433536488
                }, 
                {
                  "PWR_AttrName": "PWR_ATTR_FREQ_LIMIT_MIN", 
                  "PWR_AttrValue": 1200000, 
                  "PWR_ReturnCode": 0, 
                  "PWR_TimeNanoseconds": 80883584, 
                  "PWR_TimeSeconds": 1433536488
                }
              ], 
              "PWR_ErrorMessages": null, 
              "PWR_MajorVersion": 0, 
              "PWR_Messages": null, 
              "PWR_MinorVersion": 1, 
              "PWR_ReturnCode": 0
            } 
          }
        ]
      }

   :>json int e: Overall request status code, non-zero on partial success or
        failure

   :>json string err_msg: Human readable error message

   :>json array[object] nids: Object array containing NID specific result data, each
        element represents a single NID

   :>json int nids.nid: NID number owning the returned attribute objects

   :>json object nids.data: Remote API call result object

   :>json array[object] nids.data.PWR_Attrs: Array of attribute name value pairs

   :>json string nids.data.PWR_Attrs.PWR_AttrName: Attribute name
        copied from original request

   :>json array[int] nids.data.PWR_Attrs.PWR_AttrValue: Returned attribute
        value, in Hz

   :>json int nids.data.PWR_Attrs.PWR_ReturnCode: Attribute probe specific
        result code, non-zero on failure

   :>json int nids.data.PWR_Attrs.PWR_TimeSeconds: Elapsed seconds since
        the epoch, in UTC

   :>json int nids.data.PWR_Attrs.PWR_TimeNanoseconds: Nanosecond timestamp
        component

   :>json int nids.data.PWR_ReturnCode: Per-node attribute probe result code,
        non-zero on failure

   :>json string nids.data.PWR_ErrorMessages: Per-node attribute error
        message, or null

   :>json string nids.data.PWR_Messages: Per-node attribute info message, or
        null

   :>json int nids.data.PWR_MajorVersion: Remote API call major version number

   :>json int nids.data.PWR_MinorVersion: Remote API call minor version number

   :status 200: Network API call success
   :status 400: Bad Request
   :status 500: Internal command failure
   :status 504: Gateway Timeout

set_freq_limits
---------------
.. http:post:: /capmc/cnctl
   :synopsis: Set processor frequency limits

   Set processor operating frequency limits.

   The **set_freq_limits** call allows third-party software to constrain
   the range of frequencies a node's processor(s) may operate at. Valid
   values for minimum and maximum frequency are returned via the
   **get_freq_capabilities** call.

   This call may interact with the **aprun --p-state** switch. The ALPS
   command switch instructs the host processor to run at a fixed frequency
   for the duration of the application run. If the user requested
   performance state is not within the range specified in the
   **set_freq_limits** call, the actual performance state will be capped or
   floored such that it remains within the specified range.


    The request must **POST** a JSON object to the API server containing a
    remote API call function name, a list containing the appropriate attribute
    names, requested frequency limit values, and the protocol version. Due to
    internal bounds checking within the processor, it **is** necessary to
    specify the minimum, maximum, and a repeated minimum attribute value. 

   **Example request**:

   .. sourcecode:: http

      POST /capmc/cnctl HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 484
      Content-Type: application/json

      {
        "nids": [
          50
        ],
        "data": {
          "PWR_Function": "PWR_ObjAttrSetValues", 
          "PWR_Attrs": [
            {
              "PWR_AttrName": "PWR_ATTR_FREQ_LIMIT_MIN", 
              "PWR_AttrValue": "1400000"
            }, 
            {
              "PWR_AttrName": "PWR_ATTR_FREQ_LIMIT_MAX", 
              "PWR_AttrValue": "2500000"
            },
            {
              "PWR_AttrName": "PWR_ATTR_FREQ_LIMIT_MIN", 
              "PWR_AttrValue": "1400000"
            } 
          ], 
          "PWR_MajorVersion": 0, 
          "PWR_MinorVersion": 1
        } 
      }

   :<json array[int] nids: User specified NID list, must **not** be empty

   :<json object data: Remote API call input object

   :<json string data.PWR_Function: Remote API call function name, the caller
        must specify "PWR_ObjAttrValueCapabilities"

   :<json array[object] data.PWR_Attrs: Array of attribute name value pairs

   :<json string data.PWR_Attrs.PWR_AttrName: Attribute name, the caller
        must specify array elements for attributes named
        "PWR_ATTR_FREQ_LIMIT_MIN" and "PWR_ATTR_FREQ_LIMIT_MAX"

   :<json string data.PWR_Attrs.PWR_AttrValue: Attribute value

   :<json int data.PWR_MajorVersion: Remote API call major version number,
        the caller must specify "0"

   :<json int data.PWR_MinorVersion: Remote API call minor version number,
        the caller must specify "1"

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success", 
        "nids": [
          {
            "nid": 50,
            "data": {
              "PWR_Attrs": [
                {
                  "PWR_AttrName": "PWR_ATTR_FREQ_LIMIT_MIN", 
                  "PWR_ReturnCode": 0
                }, 
                {
                  "PWR_AttrName": "PWR_ATTR_FREQ_LIMIT_MAX", 
                  "PWR_ReturnCode": 0
                },
                {
                  "PWR_AttrName": "PWR_ATTR_FREQ_LIMIT_MIN", 
                  "PWR_ReturnCode": 0
                }
              ], 
              "PWR_ErrorMessages": null, 
              "PWR_MajorVersion": 0, 
              "PWR_Messages": null, 
              "PWR_MinorVersion": 1, 
              "PWR_ReturnCode": 0
            }
          }
        ]
      }

   The result status may indicate a non-zero return code for the first minimum
   frequency limit setting. This can happen, if for example, the previous
   maximum limit setting is less than the newly requested minimum. The
   operation may be considered successful as long as the maximum and second
   minimum frequency limit setting result is zero.

   :>json int e: Overall request status code, non-zero on partial success or
        failure

   :>json string err_msg: Human readable error message

   :>json array[object] nids: Object array containing NID specific result data, each
        element represents a single NID

   :>json int nids.nid: NID number owning the returned attribute objects

   :>json object nids.data: Remote API call result object

   :>json array[object] nids.data.PWR_Attrs: Array of attribute name value pairs

   :>json string nids.data.PWR_Attrs[].PWR_AttrName: Attribute name
        copied from original request

   :>json int nids.data.PWR_Attrs.PWR_ReturnCode: Attribute set specific
        result code, non-zero on failure

   :>json int nids.data.PWR_ReturnCode: Per-node attribute set result code,
        non-zero on failure

   :>json string nids.data.PWR_ErrorMessages: Per-node attribute error
        message, or null

   :>json string nids.data.PWR_Messages: Per-node attribute info message, or
        null

   :>json int nids.data.PWR_MajorVersion: Remote API call major version number

   :>json int nids.data.PWR_MinorVersion: Remote API call minor version number

   :status 200: Network API call success
   :status 400: Bad Request
   :status 500: Internal command failure
   :status 504: Gateway Timeout


get_sleep_state_limit_capabilities
----------------------------------
.. http:post:: /capmc/cnctl
   :synopsis: Returns sleep state limit capabilities

   Get supported processor sleep states.

   The **get_sleep_state_limit_capabilities** call informs third-party
   software about supported processor sleep states. Valid sleep
   states are returned in list form. Higher sleep state numbers
   correspond to deeper sleep states.

   The request must **POST** a JSON object to the API server containing a
   remote API call function name, a single element list containing the
   appropriate attribute name, and the protocol version.

   **Example request**:

   .. sourcecode:: http

      POST /capmc/cnctl HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 242
      Content-Type: application/json

      {
        "nids": [
          50
        ],
        "data": {
          "PWR_Function": "PWR_ObjAttrValueCapabilities", 
          "PWR_Attrs": [
            {
              "PWR_AttrName": "PWR_ATTR_CSTATE_LIMIT"
            }
          ], 
          "PWR_MajorVersion": 0,
          "PWR_MinorVersion": 1
        }
      }

   :<json array[int] nids: User specified NID list, must **not** be empty

   :<json object data: Remote API call input object

   :<json string data.PWR_Function: Remote API call function name, the caller
        must specify "PWR_ObjAttrValueCapabilities"

   :<json array[object] data.PWR_Attrs: Array of attribute name value pairs

   :<json string data.PWR_Attrs.PWR_AttrName: Attribute name, the caller must
        specify "PWR_ATTR_CSTATE_LIMIT"

   :<json int data.PWR_MajorVersion: Remote API call major version number,
        the caller must specify "0"

   :<json int data.PWR_MinorVersion: Remote API call minor version number,
        the caller must specify "1"

   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success", 
        "nids": [
          {
            "nid": 50,
            "data": {
              "PWR_Attrs": [
                {
                  "PWR_AttrName": "PWR_ATTR_CSTATE_LIMIT", 
                  "PWR_AttrValueCapabilities": [
                    0, 
                    1, 
                    2, 
                    3, 
                    4, 
                    5
                  ], 
                  "PWR_ReturnCode": 0
                }
              ], 
              "PWR_ErrorMessages": null, 
              "PWR_MajorVersion": 0, 
              "PWR_Messages": null, 
              "PWR_MinorVersion": 1, 
              "PWR_ReturnCode": 0
            }
          }
        ]
      }

   :>json int e: Overall request status code, non-zero on partial success or
        failure

   :>json string err_msg: Human readable error message

   :>json array[object] nids: Object array containing NID specific result data, each
        element represents a single NID

   :>json int nids.nid: NID number owning the returned attribute objects

   :>json object nids.data: Remote API call result object

   :>json array[object] nids.data.PWR_Attrs: Array of attribute name value pairs

   :>json string nids.data.PWR_Attrs.PWR_AttrName: Attribute name
        copied from original request

   :>json array[int] nids.data.PWR_Attrs.PWR_AttrValueCapabilities: List of
        supported sleep states

   :>json int nids.data.PWR_Attrs.PWR_ReturnCode: Attribute probe specific
        result code, non-zero on failure

   :>json int nids.data.PWR_ReturnCode: Per-node attribute probe result code,
        non-zero on failure

   :>json string nids.data.PWR_ErrorMessages: Per-node attribute error
        message, or null

   :>json string nids.data.PWR_Messages: Per-node attribute info message, or
        null

   :>json int nids.data.PWR_MajorVersion: Remote API call major version number

   :>json int nids.data.PWR_MinorVersion: Remote API call minor version number

   :status 200: Network API call success
   :status 400: Bad Request
   :status 500: Internal command failure
   :status 504: Gateway Timeout


get_sleep_state_limit
---------------------
.. http:post:: /capmc/cnctl
   :synopsis: Return state number for deepest allowable sleep state

   Get deepest allowable sleep state number.

   The **get_sleep_state_limit** call returns the state number identifying the
   deepest allowable sleep on a per-node basis. The deepest allowable sleep
   state limit may be constrained from defaults using the
   **set_sleep_state_limit** call.

   The request must **POST** a JSON object to the API server containing a
   remote API call function name, a single element list containing the
   appropriate attribute name, and the protocol version.

   **Example request**:

   .. sourcecode:: http

      POST /capmc/cnctl HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 234
      Content-Type: application/json

      {
        "nids": [
          50
        ],
        "data": {
          "PWR_Function": "PWR_ObjAttrGetValues", 
          "PWR_Attrs": [
            {
              "PWR_AttrName": "PWR_ATTR_CSTATE_LIMIT"
            }
          ], 
          "PWR_MajorVersion": 0, 
          "PWR_MinorVersion": 1
        }
      }

   :<json array[int] nids: User specified NID list, must **not** be empty

   :<json object data: Remote API call input object

   :<json string data.PWR_Function: Remote API call function name, the caller
        must specify "PWR_ObjAttrGetValues"

   :<json array[object] data.PWR_Attrs: Array of attribute name value pairs

   :<json string data.PWR_Attrs.PWR_AttrName: Attribute name, the caller must
        specify "PWR_ATTR_CSTATE_LIMIT"

   :<json int data.PWR_MajorVersion: Remote API call major version number,
        the caller must specify "0"

   :<json int data.PWR_MinorVersion: Remote API call minor version number,
        the caller must specify "1"


   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success", 
        "nids": [
          {
            "nid": 50,
            "data": {
              "PWR_Attrs": [
                {
                  "PWR_AttrName": "PWR_ATTR_CSTATE_LIMIT", 
                  "PWR_AttrValue": 5, 
                  "PWR_ReturnCode": 0, 
                  "PWR_TimeNanoseconds": 777098381, 
                  "PWR_TimeSeconds": 1433536488
                }
              ], 
              "PWR_ErrorMessages": null, 
              "PWR_MajorVersion": 0, 
              "PWR_Messages": null, 
              "PWR_MinorVersion": 1, 
              "PWR_ReturnCode": 0
            }
          }
        ]
      }
   
   :>json int e: Overall request status code, non-zero on partial success or
        failure

   :>json string err_msg: Human readable error message

   :>json array[object] nids: Object array containing NID specific result data, each
        element represents a single NID

   :>json int nids.nid: NID number owning the returned attribute objects

   :>json object nids.data: Remote API call result object

   :>json array[object] nids.data.PWR_Attrs: Array of attribute name value pairs

   :>json string nids.data.PWR_Attrs.PWR_AttrName: Attribute name
        copied from original request

   :>json array[int] nids.data.PWR_Attrs.PWR_AttrValue: Returned attribute
        value

   :>json int nids.data.PWR_Attrs.PWR_ReturnCode: Attribute probe specific
        result code, non-zero on failure

   :>json int nids.data.PWR_Attrs.PWR_TimeSeconds: Elapsed seconds since
        the epoch, in UTC

   :>json int nids.data.PWR_Attrs.PWR_TimeNanoseconds: Nanosecond timestamp
        component

   :>json int nids.data.PWR_ReturnCode: Per-node attribute probe result code,
        non-zero on failure

   :>json string nids.data.PWR_ErrorMessages: Per-node attribute error
        message, or null

   :>json string nids.data.PWR_Messages: Per-node attribute info message, or
        null

   :>json int nids.data.PWR_MajorVersion: Remote API call major version number

   :>json int nids.data.PWR_MinorVersion: Remote API call minor version number

   :status 200: Network API call success
   :status 400: Bad Request
   :status 500: Internal command failure
   :status 504: Gateway Timeout



set_sleep_state_limit
---------------------
.. http:post:: /capmc/cnctl
   :synopsis: Set the deepest sleep stat a processor may enter

   Set a node's processor(s) deepest sleep state.

   The **set_sleep_state_limit** call allows third-party software to
   constrain the deepest sleep state a node's processor(s) may enter.
   Valid values for sleep state limits are returned via the
   **get_sleep_state_limit_capabilities** call.

   Disabling sleep states, in some circumstances, can result in a slight loss
   of performance. This is due in part because idle hardware components which
   may have otherwise entered a low power state are instead forced to busy
   wait. This may cause resource contention and consume excessive power,
   subtracting from the available resources of those components performing
   useful work.

   The request must **POST** a JSON object to the API server containing a
   remote API call function name, a single element list containing the
   appropriate attribute name, requested sleep state limit value, and the
   protocol version.

   **Example request**:

   .. sourcecode:: http

      POST /capmc/cnctl HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 265
      Content-Type: application/json

      {
        "nids": [
          50
        ],
        "data": {
          "PWR_Function": "PWR_ObjAttrSetValues", 
          "PWR_Attrs": [
            {
              "PWR_AttrName": "PWR_ATTR_CSTATE_LIMIT", 
              "PWR_AttrValue": "4"
            }
          ], 
          "PWR_MajorVersion": 0, 
          "PWR_MinorVersion": 1
        }
      }

   :<json array[int] nids: User specified NID list, must **not** be empty

   :<json object data: Remote API call input object

   :<json string data.PWR_Function: Remote API call function name, the caller
        must specify "PWR_ObjAttrValueCapabilities"

   :<json array[object] data.PWR_Attrs: Array of attribute name value pairs

   :<json string data.PWR_Attrs.PWR_AttrName: Attribute name, the caller
        must specify "PWR_ATTR_CSTATE_LIMIT"

   :<json string data.PWR_Attrs.PWR_AttrValue: Attribute value

   :<json int data.PWR_MajorVersion: Remote API call major version number,
        the caller must specify "0"

   :<json int data.PWR_MinorVersion: Remote API call minor version number,
        the caller must specify "1"


   **Example response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "Success", 
        "nids": [
          {
            "nid": 50,
            "data": {
              "PWR_Attrs": [
                {
                  "PWR_AttrName": "PWR_ATTR_CSTATE_LIMIT", 
                  "PWR_ReturnCode": 0
                }
              ], 
              "PWR_ErrorMessages": null, 
              "PWR_MajorVersion": 0, 
              "PWR_Messages": null, 
              "PWR_MinorVersion": 1, 
              "PWR_ReturnCode": 0
            }
          }
        ]
      }

   :>json int e: Overall request status code, non-zero on partial success or
        failure

   :>json string err_msg: Human readable error message

   :>json array[object] nids: Object array containing NID specific result data, each
        element represents a single NID

   :>json int nids.nid: NID number owning the returned attribute objects

   :>json object nids.data: Remote API call result object

   :>json array[object] nids.data.PWR_Attrs: Array of attribute name value pairs

   :>json string nids.data.PWR_Attrs.PWR_AttrName: Attribute name
        copied from original request

   :>json int nids.data.PWR_Attrs.PWR_ReturnCode: Attribute set specific
        result code, non-zero on failure

   :>json int nids.data.PWR_ReturnCode: Per-node attribute set result code,
        non-zero on failure

   :>json string nids.data.PWR_ErrorMessages: Per-node attribute error
        message, or null

   :>json string nids.data.PWR_Messages: Per-node attribute info message, or
        null

   :>json int nids.data.PWR_MajorVersion: Remote API call major version number

   :>json int nids.data.PWR_MinorVersion: Remote API call minor version number

   :status 200: Network API call success
   :status 400: Bad Request
   :status 500: Internal command failure
   :status 504: Gateway Timeout
