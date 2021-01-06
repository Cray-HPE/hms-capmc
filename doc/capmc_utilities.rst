.. Copyright 2015,2016 Cray Inc. All Rights Reserved.

Utility Functions
==================

CAPMC contains several utility functions which may be used to query
information such as node to component name mapping, system partition
membership, or a node's role assignment. The client utility may also function
in an HTTP pass-through mode, effectively allowing a caller to specify custom
input payloads not possible using command line arguments alone. This mode of
operation offers maximum flexibility while allowing the caller to utilize
capmc's built in authorization mechanism.

get_nid_map
-----------
.. http:post:: /capmc/get_nid_map
   :synopsis: Returns NID to component name map

   Get NID to component name mapping

   Some commands, specifically commands relating to power capping
   controls, require a NID list which does not contain service nodes.
   Other use cases may require associating a geographic component name to
   NID number. In such cases, calling applications may query a node's role
   or cname using the **get_nid_map** call. The call returns a list
   objects where each element represents a single node. Each object in the
   list contains a numeric id identifying the NID number, geographic
   component name, and the operational role.

   Specify the NIDs to retreive NID number to component name mapping and
   assigned role. The syntax allows a comma-separated list of nids (eg,
   "1,4,5"), a range of nids (eg, "7-10"), or both (eg, "1,4,5,7-10"). If
   omitted, the default is all NIDs.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes a single argument which identifies a target NID
   list.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_nid_map HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 45
      Content-Type: application/json

      {
        "nids": [
          1, 
          140, 
          141
        ]
      }

   :<json array[int] nids: User specified list, or empty array for all NIDs. 

   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "", 
        "nids": [
          {
            "cname": "c0-0c2s3n1", 
            "nid": 141, 
            "role": "compute"
          }, 
          {
            "cname": "c0-0c2s3n0", 
            "nid": 140, 
            "role": "compute"
          }, 
          {
            "cname": "c0-0c0s0n1", 
            "nid": 1, 
            "role": "service"
          }
        ]
      }

   :>json int e: Error status, non-zero indicates failure
   
   :>json string err_msg: Human readable error string 

   :>json array[object] nids: Object array containing node specific mapping
        information, each element represents a single NID

   :>json int nids.nid: NID number owning the returned geographical component
        name and service role
   
   :>json string nids.cname: Geographical component name

   :>json string nids.role: Currently assigned service role, may be one of
        "compute" or "service"
   
   :status 200: Network API call success
   :status 500: Internal command failure


get_partition_map
-----------------
.. http:post:: /capmc/get_partition_map
   :synopsis: Return partition map

   Get partition map.

   Some commands, specifically commands relating to frequency and sleep
   state controls, require a NID list which does not cross system
   partition boundaries.  In such cases, calling applications may query
   partition membership using the **get_partition_map** call.  The call
   returns a list objects where each element represents a single system
   partition.  Each object in the list contains a numeric id identifying
   the partition number and the corresponding member NIDs.

   The request must **POST** an empty JSON object to the API server. This
   command takes no arguments.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_partition_map HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 2
      Content-Type: application/json

      {}

   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "", 
        "partitions": [
          {
            "partition": 0,
            "nids": [
              0, 
              3, 
              166, 
              167
            ] 
          }
        ]
      }

   :>json int e: Error status, non-zero indicates failure
   :>json string err_msg: Human readable error string 
   :>json array[object] partitions: Object array containing partition membership
        information, each element represents a single system partition
   :>json int partitions.partition: Partition number owning the returned NID
        membership list
   :>json array[int] partitions.nids: Partition member NID list

   :status 200: Network API call success
   :status 500: Internal command failure

json
----

The client script, capmc, provides a function which allows a caller to
construct and send a JSON formatted object to a user specified API handler.
This command allows a caller to utilize the capmc authorization mechanism
while constructing their own input parameter objects directly. This command is
implemented purely within the client side script.

CLI Interface
^^^^^^^^^^^^^

.. program:: json

**capmc json** 

.. option:: -r, --resource </capmc/api/path>

    Post a JSON text data structure acquired from standard input to the
    specified resource on the server. JSON text input is limited to 10MB.

.. code-block:: bash

  $ capmc json --resource=/capmc/node_status < /path/to/input/data.json

