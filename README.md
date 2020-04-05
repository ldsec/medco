# medco-connector
*medco-connector* orchestrates the MedCo query at the clinical site. It communicates with [*medco-unlynx*](https://github.com/ldsec/medco-unlynx) to execute the distributed cryptographic protocols.

## Getting started
Run the following commands to download and build the *medco-connector* module.
```shell
git clone https://github.com/ldsec/medco-connector.git
cd medco-connector/deployment/
docker-compose build
``` 

## How to use it
*medco-connector* is part of the MedCo stack. To use it, you need to have the whole MedCo stack up and running on your machine. To achieve that you can follow, for example, the [Local Development Deployment guide](https://ldsec.gitbook.io/medco-documentation/developers/local-development-deployment). 

The *medco-connector* APIs are defined using [go-swagger](https://github.com/go-swagger/go-swagger). To modify them, you must modify the `swagger/medco-connector.yml` file. To re-generate the server and client code you can run:
```shell
make swagger-gen
``` 

## How to use the medco-cli-client
*medco-connector* provides a client command-line interface to interact with the *medco-connector* APIs.

To learn how to use it, check the [CLI documentation](https://ldsec.gitbook.io/medco-documentation/system-administrators/cli).

## Source code organization

MedCo Connector uses [go-swagger](https://github.com/go-swagger/go-swagger) to generate server, client and models code.
As such, the code is organized around the Swagger definitions located in 
[swagger/medco-connector.yml](swagger/medco-connector.yml).

- *client*: client logic
- *cmd*: runnable applications
    - *cmd/medco-cli-client*: client application
    - *cmd/medco-connector-server*: server application
- *deployment*: docker image definition and docker-compose deployment
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
- *wrappers*: client library wrappers for external service (i2b2, unlynx, etc.)

## Useful information
*medco-connector* is part of the MedCo system.

You can find more information about the MedCo project [here](https://medco.epfl.ch/).

For further details, support, and contacts, you can check the [MedCo Technical Documentation](https://ldsec.gitbook.io/medco-documentation/).

## License
*medco-connector* is licensed under a End User Software License Agreement ('EULA') for non-commercial use.
If you need more information, please contact us.