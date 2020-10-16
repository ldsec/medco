# connector
The *connector* package orchestrates the MedCo query at the clinical site. 
It communicates with the [*unlynx wrapper*](../unlynx) to execute the distributed cryptographic protocols.

## Source code organization
- *client*: client logic
- *restapi*: go-swagger generated code for REST API
    - *restapi/client*: client-related generated code
    - *restapi/models*: generated models
    - *restapi/server*: server-related generated code
- *server*: server logic
    - *server/handlers*: server REST API endpoints handlers
- *swagger*: swagger REST API definitions
- *util*: utility code (configuration, security, etc.)
    - *client*: client-related utility code
    - *server*: server-related utility code
    - *common*: utility code common to client and server
- *wrappers*: client library wrappers for external service (i2b2, unlynx, etc.)

## Swagger
The *connector* APIs are defined using [go-swagger](https://github.com/go-swagger/go-swagger) by the file 
[swagger/medco-connector.yml](swagger/medco-connector.yml).
To re-generate the server, client and models code you can run from the root of the repository:
```shell
make swagger-gen
``` 
