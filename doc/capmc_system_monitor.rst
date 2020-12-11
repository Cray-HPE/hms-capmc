.. Copyright 2015,2016 Cray Inc. All Rights Reserved.

System Level Monitoring
=======================

CAPMC API calls are provided such that external agents may monitor near real
time system level power consumption and energy usage. If a time window is
specified, then historical records may be retrieved. Data may be returned in
aggregate or constituent form containing information relating to total system
or per-cabinet components, respectively. Additionally, a mechanism is provided
for a system administrator to convey intent or other operational parameters,
such as a maximum system power limit, or unreported static power overhead, to
third-parties.


get_system_parameters
---------------------
.. http:post:: /capmc/get_system_parameters
   :synopsis: Return system parameters

   Get system parameters.

   Read-only parameters such as expected worst case system power
   consumption, static power overhead, or administratively defined
   values such as a system wide power limit, maximum power ramp rate,
   and target power band may be returned via the
   **get_system_parameters** call. Returned values are used to convey
   intent between the system administrator and external agents with
   respect to target power limits and other operational parameters. The
   returned parameters are strictly informational.

   The request must **POST** an empty JSON object to the API server. This
   command takes no arguments.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_system_parameters HTTP/1.1
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
        "power_cap_target": 0, 
        "power_threshold": 0, 
        "static_power": 10700,
        "ramp_limited": false,
        "ramp_limit": 2000000,
        "power_band_min": 1000000,
        "power_band_max": 2000000
      }

   :>json int e: Error status, non-zero indicates failure
   :>json string err_msg: Human readable error string indicating failure reason
   :>json int power_cap_target: Administratively defined upper limit on system
        power
   :>json int power_threshold: System power level, which if crossed, will
        result in Cray management software emitting over power budget warnings
   :>json int static_power: Additional static system wide power overhead which
        is unreported, specified in watts

   :>json bool ramp_limited: true if out-of-band HSS power ramp rate limiting
        features are enabled
   :>json int ramp_limit: Administratively defined maximum rate of change
        (increasing or decreasing) in system wide power consumption, specified
        in watts per minute
   :>json int power_band_min: Administratively defined minimum allowable system
        power consumption, specified in watts
   :>json int power_band_max: Administratively defined maximum allowable system
        power consumption, specified in watts

   :status 200: Network API call success
   :status 500: Internal command failure

Parameters including **ramp_limited**, **ramp_limit**, **power_band_min**,
and **power_band_max** are new in SMW-8.0.

The **power_threshold** parameter is a mechanism in which a system
administrator may convey intent to an energy aware scheduler. It defines a
target, or a desired worst case system power usage. It is the responsibility
of the scheduler to enforce the nessecary energy aware scheduling policy in
order to comply with the administrators intent.

The **power_band_min** and **power_band_max** parameters allow an
administrator to convey an external constraint to a workload manager. For
example, a power utility company may state that a system should always consume
a minimum amount of power, and not exceed a maximum amount of power. These
values may be different than the system's minimum and maximum name plate rated
power values.

get_system_power
----------------
.. http:post:: /capmc/get_system_power
   :synopsis: Return system level power information

   Get system level power information.

   The **get_system_power** call returns system level power information
   including the average, minimum, and maximum values observed over a
   user specified time interval. If no arguments are given, then
   information is returned for an interval consisting of the most recent
   10 seconds.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes up to two optional arguments which identify a
   start time and window length.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_system_power HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 62
      Content-Type: application/json

      {
        "start_time": "2015-06-01 13:45:59", 
        "window_len": 30
      }


   :<json string start_time: Optional, sampling window start time
   :<json int window_len: Optional, sampling window length


   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0,
        "err_msg": "",
        "window_len": 30,
        "start_time": "2015-06-01 13:45:59", 
        "avg": 17488, 
        "max": 17661, 
        "min": 17340
      }

   :>json int e: Error status, non-zero indicates statistics are unavailable

   :>json string err_msg: Human readable error string indicating failure reason

   :>json int window_len: Window length in seconds in which the statistics
        have been computed, may be different from the requested value

   :>json string start_time: Window sample start time in
        **YYYY-MM-DD HH:MM:SS** format or symbolic constant
        **CURRENT_TIMESTAMP**, may be different from the requested value

   :>json int avg: Average system power computed over the time window

   :>json int max: Peak system power observed over the time window

   :>json int min: Min system power observed over the time window

   :status 200: Network API call success
   :status 500: Internal command failure

get_system_power_details
------------------------
.. http:post:: /capmc/get_system_power_details
   :synopsis: Returns per cabinet power information

   Get per cabinet power information.

   The **get_system_power_details** call returns per cabinet power
   information including the average, minimum, and maximum values
   observed over a user specified time interval. If no arguments are
   given, then information is returned for an interval consisting of the
   most recent 10 seconds.

   The request must **POST** a properly formatted JSON object to the API
   server. The command takes up to two optional arguments which identify a
   start time and window length.

   **Example Request**:

   .. sourcecode:: http

      POST /capmc/get_system_power_details HTTP/1.1
      Host: smw.example.com
      Accept: application/json
      Content-Length: 62
      Content-Type: application/json

      {
        "start_time": "2015-06-03 12:47:07", 
        "window_len": 30
      }


   :<json string start_time: Optional, sampling window start time
   :<json int window_len: Optional, sampling window length

   **Example Response**:

   .. sourcecode:: http

      HTTP/1.1 200 OK
      Content-Type: application/json

      {
        "e": 0, 
        "err_msg": "", 
        "window_len": 30,
        "start_time": "2015-06-03 12:47:07", 
          "cabinets": [
            {
              "avg": 17432.033333333333, 
              "max": 17730, 
              "min": 17069, 
              "x": 0, 
              "y": 0
            },
            {
              "avg": 17456.033333333333, 
              "max": 17738, 
              "min": 17060, 
              "x": 1, 
              "y": 0
            }
          ] 
      }


   :>json int e: Error status, non-zero indicates statistics are unavailable

   :>json string err_msg: Human readable error string indicating failure reason

   :>json int window_len: Window length in seconds in which the statistics
        have been computed, may be different from the requested value

   :>json string start_time: Window sample start time in
        **YYYY-MM-DD HH:MM:SS** format or symbolic constant
        **CURRENT_TIMESTAMP**, may be different from the requested value

   :>json object[] cabinets: Object array containing cabinet level power
        statistics, each element represents a single cabinet

   :>json double cabinets[].avg: Average cabinet power computed over the time window

   :>json int cabinets[].max: Peak cabinet power observed over the time window

   :>json int cabinets[].min: Min cabinet power observed over the time window

   :>json int cabinets[].x: Cabinet X coordinate, column address

   :>json int cabinets[].y: Cabinet Y coordinate, row address

   :status 200: Network API call success
   :status 500: Internal command failure

