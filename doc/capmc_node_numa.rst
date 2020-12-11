.. Copyright 2016,2017 Cray Inc. All Rights Reserved.
.. include:: <isonum.txt>

Node NUMA Capabilities & Control
================================

CAPMC NUMA controls implement a simple interface for querying and manipulating
Intel\ |copy| Xeon Phi\ |trade| processor x200 node NUMA capabilities and
configuration.  These following API calls enable external software to configure
nodes to meet various job or site configuration requirements.

All NUMA API calls, except for **get_numa_capabilities**, gather or control
staged configuration values. Staged values are applied at next boot.


get_numa_capabilities
---------------------
.. http:post:: /capmc/get_numa_capabilities
   :synopsis: Return NUMA Capabilities

   Get NUMA capabilities.

   The **get_numa_capabilities** call informs third-party software about
   the NUMA capabilities of Intel Xeon Phi processor nodes. Information
   returned includes a list objects where each object contains a Intel
   Xeon Phi processor NID number and a comma-separated list of NUMA
   configuration values that NID supports.  Information may be returned
   for a targeted set of Intel Xeon Phi processor NIDs or all Intel Xeon
   Phi processor NIDs in the system. Non-Intel Xeon Phi processor NID
   targets will be filtered out of a request. If an invalid or
   undiscovered NID is specified in the target list, an error response
   will be returned to the caller.  No such error will occur, however, if
   the request's NIDs list is empty.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_numa_capabilities HTTP/1.1
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
	If empty, the default is all NIDs.


   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0,
        "err_msg": "",
        "nids": [
          {
            "e": 0,
            "err_msg": "",
            "numa_cfg": "a2a,snc2,snc4,hemi,quad",
            "nid": 40
          }, 
          {
            "e": 0,
            "err_msg": "",
            "numa_cfg": "a2a,snc2,snc4,hemi,quad",
            "nid": 41
          }, 
          {
            "e": 0,
            "err_msg": "",
            "numa_cfg": "a2a,snc2,snc4,hemi,quad",
            "nid": 42
          }, 
          {
            "e": 0,
            "err_msg": "",
            "numa_cfg": "a2a,snc2,snc4,hemi,quad",
            "nid": 43
          }, 
          {
            "e": 0,
            "err_msg": "",
            "numa_cfg": "a2a,snc2,snc4,hemi,quad",
            "nid": 48
          }, 
          {
            "e": 0,
            "err_msg": "",
            "numa_cfg": "a2a,snc2,snc4,hemi,quad",
            "nid": 49
          }
        ]
      }


   :>json int e: Request status code, zero on success

   :>json string err_msg: Human readable error message

   :>json array[object] nids: Object array where each element represents a unique NID
        and its NUMA capability

   :>json int nids.nid: NID to which the numa_cfg applies

   :>json string nids.numa_cfg: CSV list of supported NUMA configuration
        parameters, empty when NID's status code is non-zero

   :>json int nids.e: NID status code, non-zero if NID is invalid

   :>json string nids.err_msg: Human readable error message

   :status 200: Network API call success

   :status 400: Bad request

   :status 500: Internal server error


get_numa_cfg
------------
.. http:post:: /capmc/get_numa_cfg
   :synopsis: Return NUMA configuration

   Get NUMA configuration.

   The **get_numa_cfg** call informs third-party software about the staged
   NUMA configuration of Intel Xeon Phi processor nodes, which is the NUMA
   configuration that will be applied at next boot. Information returned
   includes a list objects where each object contains a pair of
   attributes, a Intel Xeon Phi processor NID number and that NID's staged
   NUMA configuration. Information may be returned for a targeted set of
   Intel Xeon Phi processor NIDs or all Intel Xeon Phi processor NIDs in
   the system. Non-Intel Xeon Phi processor NID targets will be filtered
   out of a request, if specified. If an invalid or undiscovered NID is
   specified, an error response will be returned to the caller. No such
   error will occur, however, if the request's NIDs list is empty.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_numa_cfg HTTP/1.1
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
	If empty, the default is all NIDs.


   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0,
        "err_msg": "",
        "nids": [
          {
            "numa_cfg": "hemi", 
            "nid": 40
          }, 
          {
            "numa_cfg": "hemi", 
            "nid": 41
          }, 
          {
            "numa_cfg": "hemi", 
            "nid": 42
          }, 
          {
            "numa_cfg": "hemi", 
            "nid": 43
          }, 
          {
            "numa_cfg": "hemi", 
            "nid": 48
          }, 
          {
            "numa_cfg": "hemi", 
            "nid": 49
          }
        ]
      }


   :>json int e: Request status code, zero on success

   :>json string err_msg: Human readable error message

   :>json array[object] nids: Object array where each element represents a unique NID
        and its staged NUMA configuration

   :>json int nids.nid: NID to which the staged numa_cfg applies

   :>json string nids.numa_cfg: Staged NUMA configuration setting that will be
        applied at next boot

   :status 200: Network API call success

   :status 400: Bad request

   :status 500: Internal server error


set_numa_cfg
------------
.. http:post:: /capmc/set_numa_cfg
   :synopsis: Set NUMA configuration

   Set NUM configuration.

   The **set_numa_cfg** call allows third-party software to stage NUMA
   configuration on Intel Xeon Phi processor nodes. A node will apply the
   NUMA configuration at its next boot. Configuration can be targeted at a
   certain set of Intel Xeon Phi processor nodes or all Intel Xeon Phi
   processor nodes in a system. Non-Intel Xeon Phi processor NID targets,
   if specified, will be filtered out of a request. A simple object is
   returned to the caller indicating whether the operation succeeded or
   failed. If an invalid or undiscovered NID is specified, a more detailed
   error response will be returned to the caller.


   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies an object list
   containing objects that have two attributes, NID and NUMA configuration.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/set_numa_cfg HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 359
      Content-Type: application/json

      {
        "nids": [
          {
            "numa_cfg": "hemi", 
            "nid": 40
          }, 
          {
            "numa_cfg": "hemi", 
            "nid": 41
          }, 
          {
            "numa_cfg": "hemi", 
            "nid": 42
          }, 
          {
            "numa_cfg": "hemi", 
            "nid": 43
          }, 
          {
            "numa_cfg": "hemi", 
            "nid": 48
          }, 
          {
            "numa_cfg": "hemi", 
            "nid": 49
          }
        ]
      }

   :<json array[object] nids: User specified list of objects, where each object
        contains a pair of attributes, numa_cfg and a NID number. This list
        must not contain invalid or duplicate NID numbers. The API user should
        call **get_numa_capabilities** prior to calling **set_numa_cfg** in
        order to discover which Intel Xeon Phi processor NIDs are available and
        what their respective NUMA capabilities are. If invalid NID numbers are
        specified then an error will be returned.
    :<json int nids.nid: A NID number.

    :<json string nids.numa_cfg: Specify NUMA configuration mode. Valid values
	include: 'a2a,' 'snc2,' 'snc4,' 'hemi,' and 'quad'. The
	value 'reset' can be used to reset NUMA configuration back
	to the default value, which is 'a2a.' This is a required option.

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

   :status 400: Bad request

   :status 500: Internal server error


clr_numa_cfg
------------
.. http:post:: /capmc/clr_numa_cfg
   :synopsis: Clear NUMA configuration

   Clear NUMA configuration.

   The **clr_numa_cfg** call allows third-party software to reset the
   staged NUMA configuration on Intel Xeon Phi processor nodes to the
   default value of 'a2a'.  A node will apply the NUMA configuration at
   its next boot. Configuration can be targeted at a certain set of Intel
   Xeon Phi processor nodes or all Intel Xeon Phi processor nodes in a
   system. Non-Intel Xeon Phi processor NID targets, if specified, will be
   filtered out of a request. A simple object is returned to the caller
   indicating whether the operation succeeded or failed. If an invalid or
   undiscovered NID is specified, a more detailed error response will be
   returned to the caller.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/clr_numa_cfg HTTP/1.1
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
        "err_msg": "Success"
      }


   :>json int e: Request status code, zero on success

   :>json string err_msg: Human readable error message

   :status 200: Network API call success

   :status 400: Bad request

   :status 500: Internal server error

