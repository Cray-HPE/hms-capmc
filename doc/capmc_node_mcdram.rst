.. Copyright 2016,2017 Cray Inc. All Rights Reserved.
.. include:: <isonum.txt>

Node MCDRAM Capabilities & Control
==================================

CAPMC MCDRAM controls implement a simple interface for querying and
manipulating Intel\ |copy| Xeon Phi\ |trade| processor x200 node MCDRAM
capabilities and configuration. These following API calls enable external
software to configure nodes to meet various job or site configuration
requirements.

All MCDRAM API calls, except for **get_mcdram_capabilities**, gather or control
staged configuration values. Staged values are applied at next boot.


get_mcdram_capabilities
-----------------------
.. http:post:: /capmc/get_mcdram_capabilities
   :synopsis: Return mcdram capabilities

   Get MCDRAM capabilities

   The **get_mcdram_capabilities** call informs third-party software about
   the MCDRAM capabilities of Intel Xeon Phi processor nodes. Information
   returned includes a list objects where each object contains a Intel
   Xeon Phi processor NID number and a comma-separated list of MCDRAM
   configuration values that NID supports. Each string configuration value
   has an equivalent integer representation - both values represent the
   percentage of MCDRAM to be used as cache. Information may be returned
   for a targeted set of Intel Xeon Phi processor NIDs or all Intel Xeon
   Phi processor NIDs in the system.  Non-Intel Xeon Phi processor NID
   targets will be filtered out of a request. If an invalid or
   undiscovered NID is specified in the target list, an error response
   will be returned to the caller. No such error will occur, however, if
   the request's NIDs list is empty.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_mcdram_capabilities HTTP/1.1
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

   :<json int[] nids: User specified list, or empty array for all NIDs. This
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
            "mcdram_cfg": "flat,0,split,25,equal,50,cache,100", 
            "nid": 40
          }, 
          {
            "e": 0,
            "err_msg": "",
            "mcdram_cfg": "flat,0,split,25,equal,50,cache,100", 
            "nid": 41
          }, 
          {
            "e": 0,
            "err_msg": "",
            "mcdram_cfg": "flat,0,split,25,equal,50,cache,100", 
            "nid": 42
          }, 
          {
            "e": 0,
            "err_msg": "",
            "mcdram_cfg": "flat,0,split,25,equal,50,cache,100", 
            "nid": 43
          }, 
          {
            "e": 0,
            "err_msg": "",
            "mcdram_cfg": "flat,0,split,25,equal,50,cache,100", 
            "nid": 48
          }, 
          {
            "e": 0,
            "err_msg": "",
            "mcdram_cfg": "flat,0,split,25,equal,50,cache,100", 
            "nid": 49
          }
        ]
      }


   :>json int e: Request status code, zero on success

   :>json string err_msg: Human readable error message

   :>json object[] nids: Object array where each element represents a unique NID
        and its MCDRAM capability

   :>json int nids[].nid: NID to which the mcdram_cfg applies

   :>json string nids[].mcdram_cfg: CSV list of supported MCDRAM configuration
        parameters, empty when NID's status code is non-zero

   :>json int nids[].e: NID status code, non-zero if NID is invalid

   :>json string nids[].err_msg: Human readable error message

   :status 200: Network API call success

   :status 400: Bad request

   :status 500: Internal server error


get_mcdram_cfg
--------------
.. http:post:: /capmc/get_mcdram_cfg
   :synopsis: Return MCDRAM configuration

   Get MCDRAM configuration.

   The **get_mcdram_cfg** call informs third-party software about the
   staged MCDRAM configuration of Intel Xeon Phi processor nodes, which is
   the MCDRAM configuration that will be applied at next boot. Information
   returned includes a list objects where each object contains a Intel
   Xeon Phi processor NID number and that NIDs staged MCDRAM configuration
   and percentage, and DRAM and MCDRAM sizes. Information may be returned
   for a targeted set of Intel Xeon Phi processor NIDs or all Intel Xeon
   Phi processor NIDs in the system. Non-Intel Xeon Phi processor NID
   targets will be filtered out of a request, if specified.  If an invalid
   or undiscovered NID is specified, an error response will be returned to
   the caller. No such error will occur, however, if the request's NIDs
   list is empty.

   **Example Request**:

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   .. sourcecode:: http

      POST /capmc/get_mcdram_cfg HTTP/1.1
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
            "mcdram_cfg": "cache", 
            "mcdram_pct": 100, 
            "mcdram_size": "16384MB",
            "nid": 40,
            "dram_size": "96GB"
          }, 
          {
            "mcdram_cfg": "cache", 
            "mcdram_pct": 100, 
            "mcdram_size": "16384MB",
            "nid": 41,
            "dram_size": "96GB"
          }, 
          {
            "mcdram_cfg": "cache", 
            "mcdram_pct": 100, 
            "mcdram_size": "16384MB",
            "nid": 42,
            "dram_size": "96GB"
          }, 
          {
            "mcdram_cfg": "cache",
            "mcdram_pct": 100,
            "mcdram_size": "16384MB",
            "nid": 43,
            "dram_size": "96GB"
          }, 
          {
            "mcdram_cfg": "cache",
            "mcdram_pct": 100,
            "mcdram_size": "16384MB",
            "nid": 48,
            "dram_size": "96GB"
          }, 
          {
            "mcdram_cfg": "cache",
            "mcdram_pct": 100,
            "mcdram_size": "16384MB",
            "nid": 49,
            "dram_size": "96GB"
          }
        ]
      }


   :>json int e: Request status code, zero on success

   :>json string err_msg: Human readable error message

   :>json object[] nids: Object array where each element represents a unique NID
        and its staged MCDRAM configuration

   :>json int nids[].nid: NID to which the staged mcdram_cfg applies

   :>json string nids[].mcdram_cfg: Staged MCDRAM configuration setting that
        will be applied at next boot

   :>json int nids[].mcdram_pct: Percentage of MCDRAM that will be used as
        cache, based on the MCDRAM configuration value

   :>json string nids[].mcdram_size: NID's MCDRAM size with unit of measure

   :>json string nids[].dram_size: NID's DRAM size with unit of measure

   :status 200: Network API call success

   :status 400: Bad request

   :status 500: Internal server error


set_mcdram_cfg
--------------
.. http:post:: /capmc/set_mcdram_cfg
   :synopsis: Set MCDRAM configuration

   Set the MCDRAM configuration.

   The **set_mcdram_cfg** call allows third-party software to stage
   MCDRAM configuration on Intel Xeon Phi processor nodes. A node
   will apply the MCDRAM configuration at its next boot. Configuration
   can be targeted at a certain set of Intel Xeon Phi processor
   nodes or all Intel Xeon Phi processor nodes in a system. Non-Intel
   Xeon Phi processor NID targets, if specified, will be filtered
   out of a request. A simple object is returned to the caller
   indicating whether the operation succeeded or failed. If an
   invalid or undiscovered NID is specified, a more detailed error
   response will be returned to the caller.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies an object list
   containing objects objects that have two attributes, NID and MCDRAM
   configuration.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/set_mcdram_cfg HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 377
      Content-Type: application/json

      {
        "nids": [
          {
            "mcdram_cfg": "cache", 
            "nid": 40
          }, 
          {
            "mcdram_cfg": "cache", 
            "nid": 41
          }, 
          {
            "mcdram_cfg": "cache", 
            "nid": 42
          }, 
          {
            "mcdram_cfg": "cache", 
            "nid": 43
          }, 
          {
            "mcdram_cfg": "cache", 
            "nid": 48
          }, 
          {
            "mcdram_cfg": "cache", 
            "nid": 49
          }
        ]
      }

   :<json array[object] nids: User specified list of objects, where each object
        contains a pair of attributes, mcdram_cfg and a NID number. This list
        must not contain invalid or duplicate NID numbers. The API user should
        call **get_mcdram_capabilities** prior to calling **set_mcdram_cfg** in
        order to discover which Intel Xeon Phi processor NIDs are available and
        what their respective MCDRAM capabilities have.
    :<json int nids.nid: A NID number.  If invalid NID numbers
	are specified then an error will be returned.  If omitted,
	the default is all NIDs.
    :<json string mcdram_cfg: Specify MCDRAM configuration, which
	is the percentage of MCDRAM to be used as cache. Valid
	values include: 'flat' (or 0), 'split' (or 25), 'equal' (or
	50), and 'cache' (or 100). The value 'reset' can be used
	to reset MCDRAM configuration back to the default value,
	which is 'cache.' This is a required option.


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


clr_mcdram_cfg
--------------
.. http:post:: /capmc/clr_mcdram_cfg
   :synopsis: Clear MCDRAM configuration

   Clear MCDRAM configuration.

   The **clr_mcdram_cfg** call allows third-party software to reset
   the staged MCDRAM configuration on Intel Xeon Phi processor nodes
   to the default value of 'cache'. A node will apply the MCDRAM
   configuration at its next boot.  Configuration can be targeted
   at a certain set of Intel Xeon Phi processor nodes or all Intel
   Xeon Phi processor nodes in a system. Non-Intel Xeon Phi processor
   NID targets, if specified, will be filtered out of a request. A
   simple object is returned to the caller indicating whether the
   operation succeeded or failed. If an invalid or undiscovered NID
   is specified, a more detailed error response will be returned
   to the caller.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/clr_mcdram_cfg HTTP/1.1
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

   :>json array[int] nids: User specified list, or empty array for all NIDs. This
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

