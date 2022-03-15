# MedCo Deployment

## WebSocket tunneling of the Unlynx traffic
The schema below shows the communication flow within a MedCo node, highlighting how the WebSocket tunnel mechanism is
set up in the case of a network of 3 nodes and within node 0. The boxes are docker containers. Note that the `wstunnel`
container runs several processes, a server process that handles connections from the other nodes, and as many clients
processes as there are other nodes in the network (in this example 3-1=2 client processes). Note as well that all
incoming connections from the other nodes pass through the reverse proxy nginx.

```
|------------ unlynx -----------|   |---------- wstunnel ---------|   |--------- nginx --------|   :--- internet / other nodes ---:
|                               |   |  tunnel server              |   |  HTTPS server port 443 <--->  HTTP node 1 + 2             :
|  crypto protocols             |   |    * tunneled WS port 2003  <--->    * path /unlynx      |   :                              :
|    * server TCP port 2001     <--->    * unlynx TCP             |   |                        |   :                              :
|                               |   |  tunnel client for node 1   |   |                        |   :                              :
|    * client for node 1        <--->    * TCP port 3001          |   :                        :   :                              :
|                               |   |    * tunneled WS            <-------------------------------->  HTTP node 1                 :
|                               |   |  tunnel client for node 2   |   :                        :   :                              :
|    * client for node 2        <--->    * TCP port 3001          |   :                        :   :                              :
|                               |   |    * tunneled WS            <-------------------------------->  HTTP node 2                 :
|                               |   |-----------------------------|   :                        :   :                              :
|                               |                                     |                        |   :                              :
|                               |   |--------- connector ---------|   |                        |   :                              :
|                               |   |  REST API server            |   |                        |   :                              :
|  service API                  |   |    * HTTP REST port 1999    <--->    * path /medco       |   :                              :
|    * server TCP/WS port 2002  <--->    * unlynx service         |   |                        |   :                              :
|-------------------------------|   |-----------------------------|   |------------------------|   :------------------------------:
```

### Specific case of `dev-local-3nodes` deployment
In this deployment, this setup is still used, however the difference lies in the fact that there is a single instance
of nginx for the 3 virtual nodes. The WebSocket-tunneled unlynx traffic still goes through the nginx reverse proxy,
however a different HTTP path is used to discriminate the different nodes.
Additionally, there are as many wstunnel servers as there are nodes.
