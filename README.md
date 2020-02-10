## medco-deployment
*medco-deployment* contains the scripts to deploy MedCo in different scenarios.

A description of the available deployment profiles, along with a detailed guide on how to use them, is available [here](https://ldsec.gitbook.io/medco-documentation/system-administrator-guide/deployment).

## Source code organization
- *compose-profiles*: docker-compose files and parameters for different deployment profiles
    - *[dev-local-3nodes](https://app.gitbook.com/@ldsec/s/medco-documentation/system-administrator-guide/deployment/local-development-deployment)*: profile that deploys 3 MedCo nodes on a single machine for development purposes
    - *[test-local-3nodes](https://app.gitbook.com/@ldsec/s/medco-documentation/system-administrator-guide/deployment/local-test-deployment)*: profile that deploys 3 MedCo nodes on a single machine for test purposes
- *configuration-profiles*: configuration files for the different deployment profiles (cryptographic keys, certificates, etc.)
    - *dev-local-3nodes*: configuration files for the *dev-local-3nodes* profile
    - *test-local-3nodes*: configuration files for the *test-local-3nodes* profile
- *docker-images*: configuration files for the docker images that are used in the different deployment profiles
    - *i2b2*: configuration files for 12b2
    - *keycloak*: configuration files for keycloak
    - *nginx*: configuration files for nginx
    - *pgadmin*: configuration files for pgadmin
    - *postgresql*: configuration files for postgresql
- *resources*: additional configuration and utility files
    - *configuration*: keycloak configuration files
    - *data*: script to download the test datasets
    - *profile-generation-scripts*: scripts to generate various deployment profiles files
        - *test-network*: scripts to generate the deployment profiles files for the [Network Test Deployment](https://app.gitbook.com/@ldsec/s/medco-documentation/system-administrator-guide/deployment/network-test-deployment) profile
    - *utility-scripts*: additional utility scripts

## Useful information
*medco-deployment* is part of the MedCo system.

You can find more information about the MedCo project [here](https://medco.epfl.ch/).

For further details, support, and contacts, you can check the [MedCo Technical Documentation](https://ldsec.gitbook.io/medco-documentation/).

## License
*medco-deployment* is licensed under a End User Software License Agreement ('EULA') for non-commercial use.
If you need more information, please contact us.
