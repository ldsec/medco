[![Build Status](https://travis-ci.org/ldsec/medco.svg?branch=master)](https://travis-ci.org/ldsec/medco) 
[![Go Report Card](https://goreportcard.com/badge/github.com/ldsec/medco)](https://goreportcard.com/report/github.com/ldsec/medco) 
[![Coverage Status](https://coveralls.io/repos/github/ldsec/medco/badge.svg?branch=master)](https://coveralls.io/github/ldsec/medco?branch=master)

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
- *cmd*: binaries
    - *medco-cli-client*: REST API CLI client
    - *medco-connector-server*: REST API server
    - *medco-loader*: ETL tool
    - *medco-unlynx*: Unlynx server
- *connector*: implementation of the REST API server
- *deployments*: docker-compose files, parameters and configuration for different deployment profiles
    - *[dev-local-3nodes](https://ldsec.gitbook.io/medco-documentation/system-administrator-guide/deployment/local-development-deployment)*: profile that deploys 3 MedCo nodes on a single machine for development purposes
    - *[test-local-3nodes](https://ldsec.gitbook.io/medco-documentation/system-administrator-guide/deployment/local-test-deployment)*: profile that deploys 3 MedCo nodes on a single machine for test purposes
- *loader*: implementation of the ETL tool
- *scripts*: various utility scripts
    - *profile-generation-scripts*: scripts to generate various deployment profiles files
        - *test-network*: scripts to generate the deployment profiles files for the [Network Test Deployment](https://ldsec.gitbook.io/medco-documentation/system-administrator-guide/deployment/network-test-deployment) profile
- *test*: testing scripts
    - *test/data*: script to download the test datasets
- *unlynx*: implementation of the unlynx wrapper

## Getting started
A description of the available deployment profiles, along with a detailed guide on how to use them, is available 
[here](https://ldsec.gitbook.io/medco-documentation/system-administrator-guide/deployment).

### Building docker images
Run the following commands to build the MedCo docker images from source.
```shell
git clone https://github.com/ldsec/medco.git
cd medco/deployments/dev-local-3nodes
docker-compose -f docker-compose.yml -f docker-compose.tools.yml build
```

### Downloading docker images
Run the following commands to download the MedCo docker images.
```shell
git clone https://github.com/ldsec/medco.git
cd medco/deployments/test-local-3nodes
docker-compose -f docker-compose.yml -f docker-compose.tools.yml pull
```

### How to load data
A detailed up-to-date guide on how to use the *medco-loader* is available 
[here](https://ldsec.gitbook.io/medco-documentation/system-administrators/data-loading).

### How to use the medco-cli-client
MedCo provides a client command-line interface to interact with the *connector* APIs.

To learn how to use it, check the [CLI documentation](https://ldsec.gitbook.io/medco-documentation/system-administrators/cli).

## License
*medco-deployment* is licensed under a End User Software License Agreement ('EULA') for non-commercial use.
If you need more information, please contact us.
