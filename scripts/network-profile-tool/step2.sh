#!/bin/bash
set -Eeo pipefail
shopt -s nullglob

# ===================== input parsing ==================
function example {
    echo -e "example: $0 -nn <network name> -ni <node index> [-s <secret0,secret1,...>]"
}

function usage {
    echo -e "Wrong arguments, usage: bash $0 MANDATORY [OPTIONAL]"
}

function help {
    echo -e "MANDATORY:"
    echo -e "  -nn,   --network_name  VAL  Network name (e.g. test-network-deployment)"
    echo -e "  -ni,   --node_index    VAL  Node index (e.g. 0, 1, 2)"
    echo -e "OPTIONAL:"
    echo -e "  -s,    --secrets       VAL  Secret0,Secret1,..."
    echo -e "  -h,    --help \n"
    example
}

#Declare the number of mandatory args
margs=2

# Ensures that the number of passed args are at least equals to the declared number of mandatory args.
# It also handles the special case of the -h or --help arg.
function margs_precheck {
	if [ "$2" ] && [ "$1" -lt $margs ]; then
		if [ "$2" == "--help" ] || [ "$2" == "-h" ]; then
			help
			exit
		else
	    usage
	    help
	    exit 1
		fi
	fi
}

# Ensures that all the mandatory args are not empty
function margs_check {
	if [ $# -lt $margs ]; then
	    usage
	    help
	    exit 1 # error
	fi
}

# check if no inputs where selected
if [ $# -lt 1 ]; then
  usage
  help
  exit 1
fi
margs_precheck $# "$1"

# default values
NETWORK_NAME=
NODE_IDX=
SECRETS=

# Args while-loop
while [ "$1" != "" ];
do
   case $1 in
   -nn  | --network_name )  shift
                          NETWORK_NAME=$1
                		      ;;
   -ni  | --node_index )  shift
   						            NODE_IDX=$(printf "%03d" "$1")
			                    ;;
	 -s  | --secrets )      shift
     						          SECRETS=$1
  			                  ;;
   -h   | --help )        help
                          exit
                          ;;
   *)
                          echo "$0: illegal option $1"
                          usage
                          help
						              exit 1
                          ;;
    esac
    shift
done

# Check if all mandatory args have assigned values
margs_check "$NETWORK_NAME" "$NODE_IDX"

set -u
# convenience variables
PROFILE_NAME="network-${NETWORK_NAME}-node${NODE_IDX}"
SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MEDCO_DOCKER="ghcr.io/ldsec/medco:${MEDCO_SETUP_VER:-$(shell make --no-print-directory -C ../../ medco_version)}"
COMPOSE_FOLDER="${SCRIPT_FOLDER}/../../deployments/${PROFILE_NAME}"
CONF_FOLDER="${COMPOSE_FOLDER}/configuration"
if [[ ! -d ${CONF_FOLDER} ]] || [[ ! -d ${COMPOSE_FOLDER} ]] || [[ -f ${CONF_FOLDER}/group.toml ]]; then
    echo "The compose or configuration profile folder does not exist, or the step 2 has already been executed. Aborting."
    exit 2
fi

read -rp "### About to finalize the configuration of node ${NODE_IDX} for profile ${PROFILE_NAME}, <Enter> to continue, <Ctrl+C> to abort."

# ===================== archive extraction ==================
echo "### Extracting archives from other nodes"
pushd "${CONF_FOLDER}"
for archive in $(ls -1 srv*-public.tar.gz); do
    tar -zxvf "${archive}"
done
popd
echo "### Archives extracted"


# ===================== unlynx keys ====================
echo "### Generating group.toml and aggregate.txt files"
cat "${CONF_FOLDER}"/srv*-public.toml > "${CONF_FOLDER}/group.toml"
docker run -v "${CONF_FOLDER}:/medco-configuration" -u "$(id -u):$(id -g)" "${MEDCO_DOCKER}" medco-unlynx \
    server getAggregateKey --file "/medco-configuration/group.toml"
echo "### group.toml and aggregate.txt files generated"

echo "### Generating secrets"
if [[ -z ${SECRETS} ]]; then
    docker run -v "${CONF_FOLDER}:/medco-configuration" -u "$(id -u):$(id -g)" "${MEDCO_DOCKER}" medco-unlynx \
        server generateTaggingSecrets --file "/medco-configuration/group.toml" --nodeIndex "${NODE_IDX}"
else
    docker run -v "${CONF_FOLDER}:/medco-configuration" -u "$(id -u):$(id -g)" "${MEDCO_DOCKER}" medco-unlynx \
        server generateTaggingSecrets --file "/medco-configuration/group.toml" --nodeIndex "${NODE_IDX}" --secrets "${SECRETS}"
fi
echo "### secrets generated"


# ===================== compose profile =====================
echo "### Updating compose profile"
declare -a MEDCO_NODES_URL OIDC_JWKS_URLS OIDC_JWT_ISSUERS OIDC_CLIENT_IDS OIDC_JWT_USER_ID_CLAIMS
pushd "${CONF_FOLDER}"
ITER_IDX=0
for nodednsname in srv*-nodednsname.txt; do
  MEDCO_NODES_URL+=("https://$(<"${nodednsname}")/medco")
  OIDC_JWKS_URLS+=("https://$(<"${nodednsname}")/auth/realms/master/protocol/openid-connect/certs")
  OIDC_JWT_ISSUERS+=("https://$(<"${nodednsname}")/auth/realms/master")
  OIDC_CLIENT_IDS+=("medco")
  OIDC_JWT_USER_ID_CLAIMS+=("preferred_username")

  ITER_IDX=$((ITER_IDX+1))
done
popd

export IFS=,
cat >> "${COMPOSE_FOLDER}/.env" <<EOF
MEDCO_NODES_URL=${MEDCO_NODES_URL[*]}
OIDC_JWKS_URLS=${OIDC_JWKS_URLS[*]}
OIDC_JWT_ISSUERS=${OIDC_JWT_ISSUERS[*]}
OIDC_CLIENT_IDS=${OIDC_CLIENT_IDS[*]}
OIDC_JWT_USER_ID_CLAIMS=${OIDC_JWT_USER_ID_CLAIMS[*]}
EOF
echo "### DONE! MedCo profile ready"
