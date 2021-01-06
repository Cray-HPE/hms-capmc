.. Copyright 2015,2016 Cray Inc. All Rights Reserved.

Node Energy Reporting
=====================

CAPMC API calls are provided such that external agents may query node energy
consumption in various ways. The caller may query statistics by application
id, job id, or simply use a list of NID numbers and a caller specified time
window. Types of information returned may include aggregated energy usage on a
set of nodes, energy usage per-node, or an energy accumulator point in time
snapshot.

.. warning::

    The **get_node_energy** and **get_node_energy_stats** API calls are
    resource intensive. Depending on system size and input parameters, those
    API calls may require several minutes to complete.

get_node_energy
---------------

.. http:post:: /capmc/get_node_energy
   :synopsis: Return accumuated energy values 

   Get accumulated energy values

   Accumulated energy values for a set of nodes defined by a job (job_id),
   application (apid), list of NIDs (nids), or start and end time may be
   queried through the **get_node_energy** call. The input parameters are
   treated as conditions which are logically ANDed together. If an apid
   and start/end times are specified, then the values returned will be for
   the nodes involved in that apid during the interval specified by the
   start/end times.

   Parameters returned include the following:

   * Node count
   * Duration of the interval, in seconds
   * An array of (NID, energy) pairs

   The request must **POST** a properly formatted JSON object to the API
   server. At a minimum, the request must contain a starting and ending time
   stamp with a NID list, an ALPS application ID, or a batch scheduler ID.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_node_energy HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 122
      Content-Type: application/json

      {
        "start_time": "2015-06-03 14:07:32",
        "end_time": "2015-06-03 14:12:32", 
        "nids": [
          23, 
          24, 
          25
        ]
      }

   **Example Request**: (By apid)

   .. sourcecode:: http

      POST /capmc/get_node_energy HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 21
      Content-Type: application/json

      {
        "apid": 2977782
      }

   :<json string start_time: Optional, requested energy window sample start
	time. If omitted, the starting time of the specified **apid**
	or **jobid** is used.

   :<json string end_time: Optional, requested energy window sample end
	time.  If omitted, the ending time of the specified **apid**
	or **jobid** is used.

   :<json array[int] nids: Optional, list of NIDs to use in energy counter query.

   :<json int apid: Optional, ALPS application ID
    Return statistics applicable to the NID list and application start & end
    times for the specified ALPS id.

   :<json string job_id: Optional, batch scheduler job id
    Return statistics applicable to the NID list and application start & end
    times for the specified batch scheduler job id.


   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "", 
        "nid_count": 3, 
        "time": 300.0,
        "nodes": [
          {
            "nid": 23,
            "energy": 62623
          }, 
          {
            "nid": 24,
            "energy": 45454
          }, 
          {
            "nid": 25,
            "energy": 42870
          }
        ]
      }

   :>json int e: Error status, non-zero indicates statistics are unavailable
   
   :>json string err_msg: Human readable error string indicating failure reason
   
   :>json int nid_count: Number of nodes used in statistics query
   
   :>json double time: Window width of energy calculation, in seconds

   :>json array[object] nodes: Object array containing node level energy info, each
        element represents a single node

   :>json int nodes.nid: NID number owning the returned energy accumulation
   
   :>json int nodes.energy: Accumulated energy computed over the requested
        time interval, specified in Joules

   :status 200: Network API call success
   :status 504: Gateway Timeout
   :status 500: Internal command failure
   


get_node_energy_stats
---------------------
.. http:post:: /capmc/get_node_energy_stats
   :synopsis: Returns node energy statistics

   Get node energy statistics.

   Energy statistics for a set of nodes defined by a job (job_id),
   application (apid), list of NIDs (nodes), or start and end time may be
   queried through the **get_node_energy_stats** call. The input
   parameters are treated as conditions which are logically ANDed
   together. If an apid and start/end times are specified, then the
   statistics will be for the nodes involved in that apid during the
   interval specified by the start/end times. Both a temporal argument
   (apid, job_id, or start_time and end_time) and a component argument
   (apid, job_id, or NIDs) are required.

   Parameters returned include the following:

   * Total energy for the set
   * Average energy for nodes in the set
   * Standard deviation of energy for nodes in the set
   * An ordered pair (NID, energy) for the minimum and maximum energy consuming nodes
   * Duration of the interval, in seconds
   * Node count

   The request must **POST** a properly formatted JSON object to the API
   server. At a minimum, the request must contain a starting and ending time
   stamp with a NID list, an ALPS application ID, or a batch scheduler ID.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_node_energy_stats HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 122
      Content-Type: application/json

      {
        "start_time": "2015-06-03 14:07:32",
        "end_time": "2015-06-03 14:12:32", 
        "nids": [
          23, 
          24, 
          25
        ]
      }

   :<json string start_time: Optional, requested energy window sample start
	time.  If omitted, the starting time of the specified
	**apid** or **jobid** is used.

   :<json string end_time: Optional, requested energy window sample end
	time.  If omitted, the ending time of the specified **apid** or
	**jobid** is used.

   :<json array[int] nids: Optional, list of NIDs to use in energy counter query
    If omitted, the default is all NIDs matching other parameters.

   :<json int apid: Optional, ALPS application ID

   :<json string job_id: Optional, batch scheduler job id



   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "", 
        "time": 300.0,
        "nid_count": 3, 
        "energy_total": 150947,
        "energy_avg": 50315.666666666664, 
        "energy_std": 8766.303072307936, 
        "energy_max": [
          23, 
          62623
        ], 
        "energy_min": [
          25, 
          42870
        ]
      }

   :>json int e: Error status, non-zero indicates statistics are unavailable

   :>json string err_msg: Human readable error string indicating failure reason

   :>json double time: Window width of energy statistics calculation, in
        seconds

   :>json int energy_total: Sum of per node energy, in Joules

   :>json double energy_avg: Per node average energy, in Joules

   :>json double energy_std: Standard deviation, in Joules

   :>json array[int] energy_max: Ordered list identifying NID number followed by
        maximum observed node energy consumption

   :>json array[int] energy_min: Ordered list identifying NID number followed by
        minimum observed node energy consumption

   :status 200: Network API call success
   :status 504: Gateway Timeout
   :status 500: Internal command failure


get_node_energy_counter
-----------------------
.. http:post:: /capmc/get_node_energy_counter
   :synopsis: Return node energy counter

   Get node energy counters.

   Energy counters for a set of nodes defined by a job (job_id),
   application (apid), or list of NIDs (nids) may be queried through the
   **get_node_energy_counter** call. The parameters apid, jobid, and nids
   are treated as selectors for the set of nodes to query. If an apid or
   jobid are supplied, the running counters for each node in that aprun or
   job will be returned. If a list of NIDs is supplied, then the counters
   for the nodes corresponding to the supplied NID list will be returned.
   One and only one of apid, jobid, or NIDs list must be specified. If a
   time value is specified, then the query will retrieve the energy
   counters at or very near the specified time (if available, within one
   second). Otherwise, the most recent energy counter value will be
   returned.

   This API call returns a free running energy counter for each of the target
   NIDs. In order to be meaningful, such as when computing average power or
   energy consumed over a time interval, multiple calls must be made such
   that the caller can perform calculations based on the difference in
   returned energy counter values.

   The request must **POST** a properly formatted JSON object to the API
   server. At a minimum, the request must contain a time stamp with a NID
   list, an ALPS application ID, or a batch scheduler ID.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_node_energy_counter HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 78
      Content-Type: application/json

      {
        "time": "2015-06-03 14:07:32",
        "nids": [
          23, 
          24, 
          25
        ]
      }

   :<json string time: Optional, requested energy sample start time.
    If omitted, the energy point in time will be taken as the most recent
    available sample in the last 30 seconds on a per node basis.

   :<json array[int] nids: Optional, list of NIDs to use in energy counter
    query.  Return energy counters applicable to the NID list of
    the specified APLS id.

   :<json int apid: Optional, ALPS application ID.
    Return energy counters applicable to the NID list of the specified batch
    scheduler job id.

   :<json string job_id: Optional, batch scheduler job id

       
   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "", 
        "nid_count": 3, 
        "nodes": [
          {
            "nid": 24, 
            "energy_ctr": 14802226, 
            "time": "2015-06-03 14:07:32.886126-05"
          }, 
          {
            "nid": 23, 
            "energy_ctr": 10196418, 
            "time": "2015-06-03 14:07:32.022648-05"
          }, 
          {
            "nid": 25, 
            "energy_ctr": 13649114, 
            "time": "2015-06-03 14:07:32.886126-05"
          }
        ]
      }
   
   :>json int e: Error status, non-zero indicates statistics are unavailable

   :>json string err_msg: Human readable error string indicating failure reason

   :>json int nid_count: Number of nodes used in statistics query

   :>json object[] nodes: Object array containing node level energy info, each
        element represents a single node

   :>json int nodes[].nid: NID number owning the returned energy counter

   :>json int nodes[].energy_cntr: Point in time energy accumulator value,
        specified in Joules

   :>json string nodes[].time: Time stamp of returned energy value,
        includes fractional seconds and timezone offset

   :status 200: Network API call success
   :status 504: Gateway Timeout
   :status 500: Internal command failure

