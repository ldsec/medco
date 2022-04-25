[![Go Report Card](https://goreportcard.com/badge/github.com/ldsec/medco)](https://goreportcard.com/report/github.com/ldsec/medco) 
[![Coverage Status](https://coveralls.io/repos/github/ldsec/medco/badge.svg?branch=dev)](https://coveralls.io/github/ldsec/medco?branch=dev)

# MedCo
You can find more information about the MedCo project [here](https://medco.epfl.ch/).
For further details, support, and contacts, you can check the [MedCo Technical Documentation](https://ldsec.gitbook.io/medco-documentation/).

## Source code organization
It follows the Golang standard project layout.
- *build/package*: docker images definitions
    - *i2b2*: i2b2 docker
    - *keycloak*: keycloak docker
    - *nginx*: nginx docker
    - *pgadmin*: pgadmin docker
    - *postgresql*: postgresql docker
    - *medco*: MedCo binaries docker
    - *wstunnel*: WebSocket tunnel for unlynx communications
- *cmd*: binaries
    - *medco-cli-client*: REST API CLI client
    - *medco-connector-server*: REST API server
    - *medco-loader*: ETL tool
    - *medco-unlynx*: Unlynx server
- *connector*: implementation of the REST API server
- *[deployments](deployments/)*: deployment documentation, docker-compose files, parameters and configuration for different deployment profiles
    - *dev-local-3nodes*: profile that deploys 3 MedCo nodes on a single machine for development or test purposes
- *loader*: implementation of the ETL tool
- *scripts*: various utility scripts
    - *[network-profile-tool](scripts/network-profile-tool/)*: documentation and tool to generate network deployment profiles files
- *test*: testing scripts
    - *test/data*: script to download the test datasets
- *unlynx*: implementation of the unlynx wrapper

## Getting started
A description of the available deployment profiles, along with a detailed guide on how to use them, is available 
[here](https://ldsec.gitbook.io/medco-documentation/system-administrators/deployment).

### Build and deploy MedCo
Please refer to the [README in the deployments folder](deployments/README.md) to know how to build and deploy MedCo.

If you wish to deploy MedCo over the network accross several nodes, please refer to the 
[README in the network script folder](scripts/network-profile-tool/README.md).

### How to load data
A detailed up-to-date guide on how to use the *medco loader* is available 
[here](https://ldsec.gitbook.io/medco-documentation/system-administrators/data-loading).

### How to use the medco-cli-client
MedCo provides a client command-line interface to interact with the *connector* APIs.

To learn how to use it, check the [CLI documentation](https://ldsec.gitbook.io/medco-documentation/system-administrators/cli).

## Architectural description of a MedCo node
### WebSocket tunneling of the Unlynx traffic
The schema below shows the communication flow within a MedCo node, highlighting how the WebSocket tunnel mechanism is
set up in the case of a network of 3 nodes and within node 0. The boxes are docker containers. Note that the `wstunnel`
container runs several processes, a server process that handles connections from the other nodes, and as many clients
processes as there are other nodes in the network (in this example 3-1=2 client processes). The clients are listening to
a port in the `3000-3999` range, for the port `3XXX`, the node index is `XXX`. Note as well that all incoming
connections from the other nodes pass through the reverse proxy nginx.

```
|------------ unlynx -----------|   |---------- wstunnel ---------|   |--------- nginx --------|   :--- internet / other nodes ---:
|                               |   |  tunnel server              |   |  HTTPS server port 443 <--->  HTTP node 1 + 2             :
|  crypto protocols             |   |    * tunneled WS port 2003  <--->    * path /unlynx      |   :                              :
|    * server TCP port 2001     <--->    * unlynx TCP             |   |                        |   :                              :
|                               |   |  tunnel client for node 1   |   |                        |   :                              :
|    * client for node 1        <--->    * TCP port 3001          |   :                        :   :                              :
|                               |   |    * tunneled WS            <-------------------------------->  HTTP node 1                 :
|                               |   |  tunnel client for node 2   |   :                        :   :                              :
|    * client for node 2        <--->    * TCP port 3002          |   :                        :   :                              :
|                               |   |    * tunneled WS            <-------------------------------->  HTTP node 2                 :
|                               |   |-----------------------------|   :                        :   :                              :
|                               |                                     |                        |   :                              :
|                               |   |--------- connector ---------|   |                        |   :                              :
|                               |   |  REST API server            |   |                        |   :                              :
|  service API                  |   |    * HTTP REST port 1999    <--->    * path /medco       |   :                              :
|    * server TCP/WS port 2002  <--->    * unlynx service         |   |                        |   :                              :
|-------------------------------|   |-----------------------------|   |------------------------|   :------------------------------:
```

#### Specific case of `dev-local-3nodes` deployment
In this deployment, this setup is still used, however the difference lies in the fact that there is a single instance
of nginx for the 3 virtual nodes. The WebSocket-tunneled unlynx traffic still goes through the nginx reverse proxy,
however a different HTTP path is used to discriminate the different nodes.
Additionally, there are as many wstunnel servers as there are nodes.

## License
*medco* is licensed under an [Academic, Non-Commercial License Agreement](LICENSE).
If you need more information, please contact us.
