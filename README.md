# MedCo Connector

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
    - *util*:
- *wrappers*: client library wrappers for external service (i2b2, unlynx, etc.)
