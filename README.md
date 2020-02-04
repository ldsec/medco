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
*medco-connector* is part of the MedCo stack. To use it, you need to have the whole MedCo stack up and running on your machine. To achieve that you can follow, for example, the [Local Development Deployment guide](https://ldsec.gitbook.io/medco-documentation/system-administrator-guide/deployment/local-development-deployment). 

The *medco-connector* APIs are defined using [go-swagger](https://github.com/go-swagger/go-swagger). To modify them, you must modify the `swagger/medco-connector.yml` file. To re-generate the server and client code you can run:
```shell
make swagger-gen
``` 

## How to use the medco-cli-client

*medco-connector* provides a client command-line interface to interact with the *medco-connector* APIs.

To use it, you must first set the MEDCO_CONNECTOR_URL parameter in `deployment/docker-compose.yml` with the URL of the medco connector you want to interact with.

To show the command line manual, run:
```shell
docker-compose -f docker-compose.yml run medco-cli-client --user [USERNAME] --password [PASSWORD] --help
``` 

```shell
NAME:
   medco-cli-client - Command-line query tool for MedCo.

USAGE:
   medco-cli-client [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   query, q                                Query the MedCo network
   genomic-annotations-get-values, gval    Get genomic annotations values
   genomic-annotations-get-variants, gvar  Get genomic annotations variants
   help, h                                 Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --user value, -u value      OIDC user login
   --password value, -p value  OIDC password login
   --token value, -t value     OIDC token
   --disableTLSCheck           Disable check of TLS certificates
   --help, -h                  show help
   --version, -v               print the version
``` 

For example, to execute the `genomic-annotations-get-values` command, you could run:
```shell
docker-compose -f docker-compose.yml run medco-cli-client --user test --password test genomic-annotations-get-values hugo_gene_symbol ac
``` 

## Source code organization
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