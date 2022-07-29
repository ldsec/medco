# Network Profile Generation Tool

The provided scripts can be used to generate a MedCo deployment and configuration profile for use over network for test
or production scenarios. It is composed of the following steps that must be done from all participating nodes:
- Preliminary step: prepare the necessary information and synchronize with other nodes;
- Script execution step 1: generate own configuration and generate (or collect) keys and certificates;
- Intermediate step: share over a separate channel a generated public archive to all other nodes, and collect public archives from all other nodes;
- Script execution step 2: aggregate collected public archives and generate final configuration.

## Preliminary step
All the participating nodes must do the following prior to start the generation of the deployment profile:
- Agree to a common and unique network name which must contain only basic characters: `a-z`, `A-Z`, `0-9` and `-` (note that underscores are NOT supported), e.g. `test-network-deployment-1`;
- Agree on the participating nodes in the network and their unique index number, which must start at 0 and increase without gap in the numbering, e.g. "0, 1, 2".

Please note also that you need to have the MedCo docker images built locally. [Please refer to the deployments README to do so.](../deployments/README.md)

## Script execution step 1
Execute the script `step1.sh` with the proper arguments in order to generate part of the deployment profile.
Some examples follow.

Generate all keys and certificates and use the same address for HTTPS and unlynx:
```shell
bash step1.sh --network_name test-network-deployment --node_index 0 \
  --https_address node0.medco.com
```

Generate all keys and certificates and use different addresses for HTTPS and unlynx:
```shell
bash step1.sh --network_name test-network-deployment --node_index 0 \
  --https_address node0.medco.com --unlynx_address 192.168.57.110
```

Generate unlynx keys and provide HTTPS certificate and key:
```shell
bash step1.sh --network_name test-network-deployment --node_index 0 \
  --https_address node0.medco.com \
  --certificate ./mycert.crt --key ./mycert.key
```

Provide HTTPS certificate and key and generate unlynx keys:
```shell
bash step1.sh --network_name test-network-deployment --node_index 0 \
  --https_address node0.medco.com \
  --public_key "<unlynx_pub_key>" --secret_key "<unlynx_sec_key>"
```

Definition of all arguments:
- `--network_name` (mandatory): network name, e.g. `test-network-deployment`
- `--node_index` (mandatory): node index, e.g. `0`
- `--https_address` (mandatory): node HTTPS address, either DNS name or IP address, e.g. `test.medco.com` or `192.168.43.22`
- `--unlynx_address` (optional): unlynx address, either DNS name or IP address, if different from node HTTPS address, e.g. `test-unlynx.medco.com` or `192.168.65.55`
- `--public_key` (optional): unlynx node public key, if it is not to be generated
- `--secret_key` (optional): unlynx node private key, if it is not to be generated
- `--certificate` (optional): filepath to certificate (*.crt), if it is not to be generated
- `--key` (optional): filepath to certificate key (*.key), if it is not to be generated
- `--medco_version` (optional): version of medco to be used

## Intermediate step
### Share public archive with all other nodes
During the execution of step 1, a public archive named like `srvXXX-public.tar.gz` has been generated in the
`configuration` folder of the deployment profile. This file must be shared to all the other nodes on a separate channel,
e.g. over email.

### Collect public archive from all other nodes
The archives (named like `srvXXX-public.tar.gz`) must be collected from all other nodes and put in the `configuration`
folder of the deployment profile before proceeding to the next step.

## Script execution step 2
Execute the script `step2.sh` with the proper arguments in order to finalize the deployment profile.
Some examples follow.

Aggregate public archives from all nodes and generate unlynx DDT secrets:
```shell
bash step2.sh --network_name test-network-deployment --node_index 0 --nb_nodes 3
```

Aggregate public archives from all nodes and provide unlynx DDT secrets:
```shell
bash step2.sh --network_name test-network-deployment --node_index 0 --nb_nodes 3 \
  --secrets "<secret0>,<secret1>,<secret2>"
```

Definition of all arguments:
- `--network_name` (mandatory): network name, e.g. `test-network-deployment`
- `--node_index` (mandatory): node index, e.g. `0`
- `--nb_nodes` (mandatory): total number of nodes in the network, e.g. `3`
- `--secrets` (optional): unlynx DDT secrets, if they are not to be generated, e.g. `<secret0>,<secret1>,<secret2>`