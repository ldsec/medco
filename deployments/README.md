# Pre-requisites
## Hardware
- Network Bandwidth: >100 Mbps (ideal), >10 Mbps (minimum), symmetrical
- Ports Opening and IP Restrictions: see Network Architecture
- Hardware
  - CPU: 8 cores (ideal), 4 cores (minimum)
  - RAM: >16 GB (ideal), >8GB (minimum)
  - Storage: dependent on data loaded, >100GB

## Software
- OS: Any flavor of Linux, physical or virtualized (tested with Ubuntu 16.04, 18.04, Fedora 29-33)
- Git
- OpenSSL
- Docker version >= 18.09.1, correctly configured:
  - to be used as a normal user
  - range of IPs in virtual network used by MedCo not conflicting with others on your system 
- Docker-Compose version >= 1.23.2

# Stack Deployment
The following steps assumed that you already retrieved the MedCo repository.

They also assume the use of the following variables that you must adapt to your own situation:
```shell
export MEDCO_REPO=/opt/medco
export MEDCO_DEPLOYMENT_PROFILE=dev-local-3nodes
```

## Load environment variables
In order to proceed with the build and/or deployment it is mandatory to load some environment variables that configure
the versions used. You can do so by sourcing a script in your shell session, *beware that you will need to repeat this
action everytime you use a new one*.

To do so, run *from the root of the repository* the following command, after replacing `<version>` with:
- `dev` for a local development environment where you build your own docker images
- leave empty for generating a default version number based on the current git HEAD
- any version, which if it exists on the container registry can be pulled, otherwise it can be built locally, 
  e.g. `v3.0.0`; if you are deploying MedCo in production this is the recommended approach

```shell
cd "${MEDCO_REPO}"
source scripts/versions.sh <version>
```

## Retrieve the Docker images by ...
### ... building them
If you are not pulling the pre-built docker images, you must build the Docker images with:
```shell
cd "${MEDCO_REPO}/deployments/${MEDCO_DEPLOYMENT_PROFILE}"
make build
```

Note the directory change.  
The images will be tagged with the versions set in your shell session during the previous step.

### ... downloading them
Note that this will only work if the versions selected at the previous steps match pre-built versions available in the
container registry!
```shell
cd "${MEDCO_REPO}/deployments/${MEDCO_DEPLOYMENT_PROFILE}"
make pull
```

Note the directory change.

## Customize pre-deployment configuration
In the root of the deployment profile you are deploying, edit the `.env` file to change:
- the default passwords;
- if needed the different log levels.

Note that some of those setting must be changed before the first deployment, because if they are not, changing them
afterward will have no effect. 

## Deploy the MedCo stack ...
### ... with the `dev-local-3nodes` profile
In the `dev-local-3nodes` folder, edit the `.env` file to change the `MEDCO_NODE_HOST` variable to the public address of
the machine your are deploying on. Note that you do not need to change it (i.e. leave it to `localhost`) if you are
deploying on your own machine.

Run the following command:
```shell
cd "${MEDCO_REPO}/deployments/${MEDCO_DEPLOYMENT_PROFILE}"
make up
```

Important: you must not interrupt the first startup! If you do the database will end up corrupted. In order to assess
if it is safe to stop the deployment, check out the logs of the `i2b2` container to see if it is still loading data or
not.
If you happen to do so, you must reset the deployment (see below) before starting it again.

### ... with a `network` profile
You must first generate a network profile (along with the other participating nodes) using the provided tool
in [../scripts/network-profile-tool](../scripts/network-profile-tool). Then run the following command from the *root of
the deployment profile you generated*:
```shell
cd "${MEDCO_REPO}/deployments/${MEDCO_DEPLOYMENT_PROFILE}"
make up
```

## Post-deployment steps
### Keycloak configuration
You must mandatorily generate new keys for the Keycloak realm (HMAC, AES, RSA).
Additionally, if it is a production deployment, you must change the passwords of all the keycloak users.

After you did so, you must restart the deployment:
```shell
cd "${MEDCO_REPO}/deployments/${MEDCO_DEPLOYMENT_PROFILE}"
make stop
make up
```

### Accept self-signed certificates in your browser
If some nodes of network are using the default self-signed certificates, or if you deployed the `dev-local-3nodes`
profile, you will need to visit in your browser the web page of the nodes in question in order for your browser to accept
those certificates. If you don't, the queries made by your browser to all of the nodes of the network will silently fail
in the background.

### Load test data
First, start by downloading the test data. *From the root of the repository*, run the following command:
```shell
cd "${MEDCO_REPO}"
make download_test_data
```

Then, *from the root of the deployment profile you are using*, run the following command:
```shell
cd "${MEDCO_REPO}/deployments/${MEDCO_DEPLOYMENT_PROFILE}"
make load_test_data
```

### Stop the deployment
In order to stop the deployment you can:
```shell
cd "${MEDCO_REPO}/deployments/${MEDCO_DEPLOYMENT_PROFILE}"
make stop
```

In order to stop the deployment, and delete the docker containers and virtual networks:
```shell
cd "${MEDCO_REPO}/deployments/${MEDCO_DEPLOYMENT_PROFILE}"
make down
```

### Handling a reverse proxy (RP-WAF) in front of the deployment
If you happen to have a RP-WAF in front of the MedCo node, it needs to be configured to:
- accept the self-signed certificate of the MedCo node (if another one is not set up);
- let through WebSocket traffic, at least on the path `/unlynx`.

## In case of issue: reset deployment
In case of issue, it is recommended to start the deployment from scratch.

Identify the named volume name:
```shell
docker volume ls
```

Then bring down the deployment and delete the named volume containing the database:
```shell
cd "${MEDCO_REPO}/deployments/${MEDCO_DEPLOYMENT_PROFILE}"
make down
docker volume rm ${MEDCO_DEPLOYMENT_PROFILE}_medcodb
```
