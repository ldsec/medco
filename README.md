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
- *deployments*: docker-compose files, parameters and configuration for different deployment profiles
    - *[dev-local-3nodes](https://ldsec.gitbook.io/medco-documentation/developers/local-development-deployment)*: profile that deploys 3 MedCo nodes on a single machine for development or test purposes
- *loader*: implementation of the ETL tool
- *scripts*: various utility scripts
    - *network-profile-tool*: scripts to generate the deployment profiles files for the [Network Deployment](https://ldsec.gitbook.io/medco-documentation/system-administrators/deployment/network-deployment) profile
- *test*: testing scripts
    - *test/data*: script to download the test datasets
- *unlynx*: implementation of the unlynx wrapper

## Getting started
A description of the available deployment profiles, along with a detailed guide on how to use them, is available 
[here](https://ldsec.gitbook.io/medco-documentation/system-administrators/deployment).

### Retrieve repository and prepare environment
Run the following commands to clone the MedCo repository and to select a version of MedCo.
See the [deployments README](deployments/) for more information on selecting the version. 

```shell
git clone https://github.com/CHUV-DS/medco.git
cd medco/deployments/dev-local-3nodes
source ../../scripts/versions.sh
```

### Build docker images
Run the following commands to build the MedCo docker images from source.
```shell
make build
```

### Download docker images
Run the following commands to download the MedCo docker images.
```shell
make pull
```

### How to load data
A detailed up-to-date guide on how to use the *medco loader* is available 
[here](https://ldsec.gitbook.io/medco-documentation/system-administrators/data-loading).

### How to use the medco-cli-client
MedCo provides a client command-line interface to interact with the *connector* APIs.

To learn how to use it, check the [CLI documentation](https://ldsec.gitbook.io/medco-documentation/system-administrators/cli).

## License
*medco* is licensed under an [Academic, Non-Commercial License Agreement](LICENSE).
If you need more information, please contact us.
