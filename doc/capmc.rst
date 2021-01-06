.. Copyright 2015 Cray Inc. All Rights Reserved.

Capmc CLI
=========

The :index:`capmc` utility provides remote monitoring and control capabilities
to agents running externally to the Cray System Management Workstation
(:index:`SMW`). The capmc client utility accepts command line arguments from
the caller and submits an appropriately formatted request via HTTPS to the
monitoring and control service providers running on the SMW. Command results
are supplied as JSON-formatted text to standard output. Remote access is
restricted by means of public-key authorization established by the site
administrator.

In order to use the capmc :index:`API` calls, also referred to as *applets*,
caller must first load the capmc module:

.. code-block:: bash

    $ module load capmc

Alternatively, the caller can invoke capmc by specifying the absolute path
(determined on installation). The default path is:
**/opt/cray/capmc/default/bin/capmc**

Installation
------------

The capmc utility is packaged in a noarch :index:`RPM` format. Its only
dependency is the Python Standard Library with SSL support enabled, version
2.6 or 2.7. On RPM based systems, the package may be installed with the native
package manager.

.. code-block:: bash

    $ rpm -ivh cray-capmc-1.0-1.0000.36550.40.1.noarch.rpm 
    Preparing...      ########################################### [100%]
       1:cray-capmc   ########################################### [100%]

Alternately, capmc may be installed on systems that do not support the RPM
format using the rpm2cpio utility. The following example extracts the package 
contents into the current working directory.

.. code-block:: bash

    $ rpm2cpio cray-capmc-1.0-1.0000.36550.40.1.noarch.rpm | cpio -idmv 
    ./etc/opt/cray
    ./etc/opt/cray/capmc
    ./opt/cray
    ./opt/cray/capmc
    ./opt/cray/capmc/1.0-1.0000.36550.40.1
    ./opt/cray/capmc/1.0-1.0000.36550.40.1/bin
    ./opt/cray/capmc/1.0-1.0000.36550.40.1/bin/capmc
    ...

Return Values
-------------
The capmc utility returns 0 to the shell if the request was successfully
submitted to the platform monitoring and control service. Individual command
status is indicated in the JSON-encoded output data 'e' member. 'e' is always
present in the response payloads and will be non-zero to indicate a command
specific error code if the requested operation failed. Callers of capmc should
always evaluate result status returned in the message payload.


Credentials
-----------

Capmc utilizes an :index:`X.509` client certificate for authorization. The
signed client certificate and private key will be supplied by the system
administrator. Credentials may be supplied through environment variables or a
configuration file. 


Environment Variables
^^^^^^^^^^^^^^^^^^^^^

The required environment variables are listed below. Environment variables
will override configuration file parameters. Environment variables which must
be set to utilize client certificate authorization are as follows: 

.. envvar:: OS_KEY

    Specify the absolute path in the local filesystem where the X.509 client
    certificate key is installed. If a pass phrase is set on the key, the
    user will be prompted when required.

.. envvar:: OS_CERT

    Specify the absolute path in the local filesystem where the X.509 client
    certificate is installed.

.. envvar:: OS_CACERT

    Specify the absolute path in the local filesystem where the X.509
    Certificate Authority (CA) certificate is located.

.. envvar:: OS_SERVICE_URL

    Specify the URL where the application service is listening. This must
    include the fully qualified domain name (FQDN) of the SMW, protocol, and
    port number. The default port number is 8443.

.. ifconfig:: version in ('dev')

    Alternatly, capmc may integrate with the OpenStack's Keystone, an
    Identity Management Service. Keystone users may be granted permission to
    use capmc by granting the role 'capmc.' Additionally, the capmc remote
    service provider called 'xtremote' must have the appropriate publicURL
    configured in the Keystone service catalog. 

    .. envvar:: OS_USERNAME

        Specify the username with permission to invoke capmc.

    .. envvar:: OS_PASSWORD

        Specify the password associated with the user.

    .. envvar:: OS_TENANT_NAME

        Specify the project name configured by the site administrator. This
        will likely be "cray". 

    .. envvar:: OS_AUTH_URL

        Specify the URL where the authorization service is listening. This
        must include the fully qualified domain name (FQDN) of the SMW,
        protocol, and port number. The default port number is 35357.

    .. envvar:: OS_CACERT

        Specify the absolute path in the local filesystem where the X.509
        Certificate Authority (CA) certificate is located.


Configuration File
^^^^^^^^^^^^^^^^^^

When capmc is utilized for autonomous machine to machine communication, it is
advisable to define the certificate paths and service url in a configuration
file instead of the callers environment.  Capmc reads its configuration file
from the following location: **/etc/opt/cray/capmc/capmc.json**

X.509 client certificate configuration syntax:

.. code-block:: javascript

   {
      "os_key":         "/etc/opt/cray/capmc/capmc-client.key",
      "os_cert":        "/etc/opt/cray/capmc/capmc-client.pem",
      "os_cacert":      "/etc/opt/cray/capmc/capmc-cacert.pem",
      "os_service_url": "https://smw.example.com:8443"
   }

.. ifconfig:: version in ('dev')

    Alternate example, Keystone client configuration syntax:

    .. code-block:: javascript

       {
        "os_username":    "<foo-user>",
        "os_password":    "<foo-password>",
        "os_tenant_name": "cray",
        "os_auth_url":    "https://smw.example.com:35357/v2.0",
        "os_cacert":      "/etc/opt/cray/capmc/capmc-cacert.pem"
       }

.. note::
    Capmc validates all X.509 certificates. It is assumed that all host
    certificates used within HTTPS API servers such as Xtremoted, or X.509
    client authorization are signed by the same certificate authority.

.. warning::
    The configuration and X.509 certificate files must have appropriate
    filesystem permissions. Users with read access to the global configuration
    file and the client certificate will be granted permission to utilize all
    capmc functionality.

