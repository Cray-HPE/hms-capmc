.. Copyright 2016,2017 Cray Inc. All Rights Reserved.
.. include:: <isonum.txt>

Node SSD Control & Diagnostics
==============================

CAPMC SSD controls implement a simple interface for querying and manipulating
Intel\ |copy| Xeon Phi\ |trade| processor x200 node on-board SSDs. These
following API calls enable external software to configure nodes to meet various
job or site configuration requirements. These calls are limited to Intel Xeon
Phi processor nodes.

All Node SSD Control API calls gather or control staged configuration values.
Staged values are applied at next boot.

CAPMC SSD diagnostics implement an interface for querying various inventory and
diagnositic information from SSDs on a system. SSD diagnostics can target any
node in a system, not just Intel Xeon Phi processor nodes.


get_ssd_enable
--------------
.. http:post:: /capmc/get_ssd_enable
   :synopsis: Return state of SSDs

   Get the state of SSDs on Intel Xeon Phi processor nodes.

   SSDs on Intel Xeon Phi processor nodes can have a state of enabled
   or disabled.  The **get_ssd_enable** call informs third-party
   software about the staged state of SSDs on Intel Xeon Phi processor
   nodes, which is the SSD state that will be applied at next boot.
   Information returned includes a list objects where each object
   contains a pair of attributes, a Intel Xeon Phi processor NID
   number and the stated state of that NID's SSD.  Information may
   be returned for a targeted set of Intel Xeon Phi processor NIDs
   or all Intel Xeon Phi processor NIDs in the system. Non-Intel
   Xeon Phi processor NID targets will be filtered out of a request,
   if specified. If an invalid or undiscovered NID is specified,
   an error response will be returned to the caller. No such error
   will occur, however, if the request's NIDs list is empty.



   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_ssd_enable HTTP/1.1
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
            "ssd_enable": 1,
            "nid": 40
          }, 
          {
            "ssd_enable": 1,
            "nid": 41
          }, 
          {
            "ssd_enable": 1,
            "nid": 42
          }, 
          {
            "ssd_enable": 1,
            "nid": 43
          }, 
          {
            "ssd_enable": 1,
            "nid": 48
          }, 
          {
            "ssd_enable": 1,
            "nid": 49
          }
        ]
      }


   :>json int e: Request status code, zero on success

   :>json string err_msg: Human readable error message

   :>json array[object] nids: Object array where each element represents a unique NID
        and its staged SSD state

   :>json int nids.nid: NID to which the staged ssd_enable setting applies

   :>json int nids.ssd_enable: Staged SSD state that will be applied at next
        boot

   :status 200: Network API call success

   :status 400: Bad request

   :status 500: Internal server error


set_ssd_enable
--------------
.. http:post:: /capmc/set_ssd_enable
   :synopsis: Set SSD state enable

   Set the SSD state to enable.

   The **set_ssd_enable** call allows third-party software to stage SSD
   state on Intel Xeon Phi processor nodes. A node will apply the SSD
   state value at its next boot.  Configuration can be targeted at a
   certain set of Intel Xeon Phi processor nodes or all Intel Xeon Phi
   processor nodes in a system. Non-Intel Xeon Phi processor NID targets,
   if specified, will be filtered out of a request.  A simple object is
   returned to the caller indicating whether the operation succeeded or
   failed. If an invalid or undiscovered NID is specified, a more detailed
   error response will be returned to the caller.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies an object list
   containing objects that have two attributes, NID and SSD state.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/set_ssd_enable HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 341
      Content-Type: application/json

      {
        "nids": [
          {
            "ssd_enable": 1,
            "nid": 40
          }, 
          {
            "ssd_enable": 1,
            "nid": 41
          }, 
          {
            "ssd_enable": 1,
            "nid": 42
          }, 
          {
            "ssd_enable": 1,
            "nid": 43
          }, 
          {
            "ssd_enable": 1,
            "nid": 48
          }, 
          {
            "ssd_enable": 1,
            "nid": 49
          }
        ]
      }

   :<json array[object] nids: User specified list of objects, where each object
        contains a pair of attributes, ssd_enable and a NID number. This list
        must not contain invalid or duplicate NID numbers. If invalid NID
	numbers are specified then an error will be returned.  If
	empty, the default is all NIDs.
   :<json int nid: A NID number.
   :<json int ssd_enable: Specify SSD state mode. Valid values
	include: 'enable,' 'disable,' and 'reset.' The value 'reset'
	can be used to reset SSD configuration back to the default
	state value, which is 'enable.' This is a required option.

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


clr_ssd_enable
--------------
.. http:post:: /capmc/clr_ssd_enable
   :synopsis: Reset staged SSD state

   Clear (reset) the staged SSD state.

   The **clr_ssd_enable** call allows third-party software to reset the
   staged SSD state on Intel Xeon Phi processor nodes to the default value
   of '1' (which is 'enable'). A node will apply the SSD state at its next
   boot. Configuration can be targeted at a certain set of Intel Xeon Phi
   processor nodes or all Intel Xeon Phi processor nodes in a system.
   Non-Intel Xeon Phi processor NID targets, if specified, will be
   filtered out of a request. A simple object is returned to the caller
   indicating whether the operation succeeded or failed. If an invalid or
   undiscovered NID is specified, a more detailed error response will be
   returned to the caller.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a NID list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/clr_ssd_enable HTTP/1.1
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

   :<json array[int] nids: User specified list, or empty array for
	all NIDs. This list must not contain invalid or duplicate
	NID numbers. If invalid NID numbers are specified then an
	error will be returned.


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


get_ssds
--------
.. http:post:: /capmc/get_ssds
   :synopsis: Return SSD inventory

   Get SSD inventory information.

   The **get_ssds** call provides third-party software with SSD inventory
   information. Information returned includes an object with zero or more
   node cname attributes that each map to a list of one or more SSD
   inventory objects.  Each inventory object can include the SSD's PCI
   bus, device and function, model number, serial number, and system and
   sub-system device and vendor IDs.  Information may be returned for a
   targeted set of NIDs or all NIDs in the system. If the targeted NIDs do
   not have SSDs attached, the response object will only contain a success
   or failure code and message for the API call itself. If an invalid NID
   is specified, an error response will be returned to the caller.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes no arguments (ie, an empty object) or a single
   argument, which is either a list of NID numbers or a list of node cname
   strings.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_ssd_enable HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 27
      Content-Type: application/json

      {
        "nids": [
          104
        ]
      }

   :<json array[int] nids: User specified list, or empty array for all NIDs. This
        list must not contain invalid or duplicate NID numbers. If invalid
        NID numbers are specified then an error will be returned.

   :<json array[string] cname: User specified list, or empty array for all nodes.
        This list must not contain invalid cnames or cnames referencing objects
        that are not nodes. If invalid cnames are specified then an error will
        be returned.

   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "err_msg": "Success",
        "e": 0,
        "c0-0c1s10n0": [
          {
            "sub_id": "144da801",
            "func": 0,
            "serial_number": "S0ABCDAH000000",
            "device": 0,
            "bus": 2,
            "ssd_id": "144da802",
            "model_number": "SAMSUNG MZVPV128HDGM-00000",
            "nid": 104,
            "size": 128000000000
          }
        ]
      }


   :>json int e: Request status code, zero on success

   :>json string err_msg: Human readable error message

   :>json array[object] cname: A list of objects containing SSD inventory
        information, one for each attached SSD. Note the label of this
        attribute will not be "cname" in the literal sense, but rather it will
        be the node's cname (eg, "c0-0c1s10n0")

   :>json int cname.bus: PCI bus to which the SSD is attached

   :>json int cname.device: SSD's PCI device number on the bus

   :>json int cname.func: SSD's PCI function identifier

   :>json string cname.model_number: SSD's model number

   :>json int cname.nid: NID number to which this SSD inventory
        information belongs

   :>json string cname.serial_number: SSD's serial number

   :>json int cname.size: Size of the SSD, not its remaining capacity

   :>json string cname.ssd_id: The SSD's PCI system device and vendor
        identifier

   :>json string cname.sub_id: The SSD's PCI sub-system device and vendor
        identifer

   :status 200: Network API call success

   :status 400: Bad request

   :status 500: Internal server error


get_ssd_diags
-------------
.. http:post:: /capmc/get_ssd_diags
   :synopsis: Return SSD diagnostic information

   Get SSD diagnostic information.

   The **get_ssd_diags** call provides third-party software with SSD
   diagnostic information. Diagnostic information is pulled from a
   database on the SMW.  Information returned includes an object that has
   an attribute with a list of SSD diagnostic objects as its value. Each
   SSD diagnostic object includes information such as a timestamp when the
   data was recorded in the database, life remaining, firmware version,
   serial_number, manufacturer's identifer, part identifier, PCI
   coordinates (bus, device, function), SSD size (not remaining capacity),
   and percentage used. Information may be returned for a targeted set of
   NIDs or all NIDs in the system. If there is no SSD diagnostic data in
   the databse for a targeted SSD, the returned information will have an
   empty list of SSD diagnostic objects. If an invalid NID is specified,
   an error response will be returned to the caller.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes no arguments (ie, an empty object) or a single
   argument, which is either a list of NID numbers or a list of node cname
   strings.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_ssd_diags HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 27
      Content-Type: application/json

      {
        "nids": [
          104
        ]
      }

   :<json array[int] nids: User specified list, or empty array for all NIDs. This
        list must not contain invalid or duplicate NID numbers. If invalid
        NID numbers are specified then an error will be returned.

   :<json array[string] cname: User specified list, or empty array for all nodes.
        This list must not contain invalid cnames or cnames referencing objects
        that are not nodes. If invalid cnames are specified then an error will
        be returned.

   **Example Response**:

   .. sourcecode::diags
      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0,
        "err_msg": "Success",
        "ssd_diags": [
          {
            "life_remaining": 100.0,
            "firmware": "8DV10170",
            "ts": "2016-06-23 09:55:25-06:00",
            "nid": 18,
            "serial_num": "CXA0000000000000-1",
            "cname": "c0-0c0s4n2",
            "manu_id": "8086",
            "percent_used": 0.0,
            "part_id": "8718",
            "comp_ord": "03:00:0",
            "size": 2000
          },
          {
            "life_remaining": 100.0,
            "firmware": "8DV10170",
            "ts": "2016-06-23 09:55:25-06:00",
            "nid": 18,
            "serial_num": "CXA0000000000000-2",
            "cname": "c0-0c0s4n2",
            "manu_id": "8086",
            "percent_used": 0.0,
            "part_id": "8718",
            "comp_ord": "03:00:0",
            "size": 2000
          }
        ]
      }


   :>json int e: Request status code, zero on success

   :>json string err_msg: Human readable error message

   :>json array[object] ssd_diags: A list of objects containing SSD diagnostic
        information

   :>json string ssd_diags.cname: Cname of node to which this SSD diagnostic
        information belongs

   :>json string ssd_diags.comp_ord: The PCI bus, device, function coordinates

   :>json string ssd_diags.firmware: Firmware version flashed on SSD

   :>json float ssd_diags.life_remaining: SSD's remaining life

   :>json string ssd_diags.manu_id: Manufacturer's identifer

   :>json int ssd_diags.nid: NID number to which this SSD diagnostic
        information belongs

   :>json string ssd_diags.part_id: Manufacturer's part identifier

   :>json float ssd_diags.percent_used: Percentage of formatted size used

   :>json string ssd_diags.serial_number: SSD's serial number

   :>json int ssd_diags.size: SSD's size in GB, not its remaining
        capacity

   :>json string ssd_diags.ts: Timestamp when data was written to the
        database

   :status 200: Network API call success

   :status 400: Bad request

   :status 500: Internal server error

