#!/bin/bash
set -Eeo pipefail
shopt -s nullglob

source common.sh

# ===================== input parsing ==================
function example {
    echo -e "example: $0 -nn <network name> -ni <node index> -nb <total number of nodes> [-s <secret0,secret1,...>]"
}

function help {
    echo -e "MANDATORY:"
    echo -e "  -nn,   --network_name  VAL  Network name (e.g. test-network-deployment)"
    echo -e "  -ni,   --node_index    VAL  Node index (e.g. 0, 1, 2)"
    echo -e "  -nb,   --nb_nodes      VAL  Total number of nodes in the network (e.g. 3)"
    echo -e "OPTIONAL:"
    echo -e "  -s,    --secrets       VAL  Unlynx DDT secrets, if they are not to be generated (e.g. <secret0>,<secret1>,<secret2>)"
    echo -e "  -h,    --help \n"
    example
}

margs=3 # number of mandatory args
margs_precheck $# "$1" $margs

# default values
NETWORK_NAME=
NODE_IDX=
NB_NODES=
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
   -nb  | --nb_nodes )    shift
   						            NB_NODES=$(printf "%03d" "$1")
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
margs_check $margs "$NETWORK_NAME" "$NB_NODES" "$NODE_IDX"

set -u
# generate convenience variables
export_variables "$NETWORK_NAME" "$NODE_IDX"
if [[ ! -d ${CONF_FOLDER} ]] || [[ ! -d ${COMPOSE_FOLDER} ]] || [[ -f ${CONF_FOLDER}/group.toml ]]; then
    echo "The compose or configuration profile folder does not exist, or the step 2 has already been executed. Aborting."
    exit 2
fi

read -rp "### About to finalize the configuration of node ${NODE_IDX} for profile ${PROFILE_NAME}, <Enter> to continue, <Ctrl+C> to abort."
dependency_check

# ===================== archive extraction ==================
pushd "${CONF_FOLDER}"

echo "### Checking number of nodes public archives"
find . -maxdepth 1 -name 'srv*-public.tar.gz'
nb_pub_archives=$(find . -maxdepth 1 -name 'srv*-public.tar.gz' | wc -l)
if [[ "$NB_NODES" -ne "$nb_pub_archives" ]]; then
  echo "ERROR: Found ${nb_pub_archives} nodes public archives, which must be equal to ${NB_NODES} but is not, exiting."
  exit 2
fi
echo "### Found correct number of nodes public archives"

echo "### Extracting archives from other nodes"
for archive in srv*-public.tar.gz
do
  [[ -e "$archive" ]] # if no archive
  tar -zxvf "${archive}"
done
echo "### Archives extracted"

popd

# ===================== unlynx keys ====================
echo "### Generating group.toml and aggregate.txt files"
cat "${CONF_FOLDER}"/srv*-public.toml > "${CONF_FOLDER}/group.toml"
"${MEDCO_BIN[@]}" medco-unlynx server getAggregateKey --file "/medco-configuration/group.toml"
echo "### group.toml and aggregate.txt files generated"

unlynx_secrets_gen_args=(
  medco-unlynx server generateTaggingSecrets
  --file "/medco-configuration/group.toml"
  --nodeIndex "${NODE_IDX}"
)

echo "### Generating unlynx secrets"
if [[ -n "$SECRETS" ]]; then
  echo "### Using pre-generated unlynx secrets"
  unlynx_secrets_gen_args=("${unlynx_secrets_gen_args[@]}" --secrets "${SECRETS}")
fi

"${MEDCO_BIN[@]}" "${unlynx_secrets_gen_args[@]}"
echo "### unlynx secrets generated"


# ===================== compose profile =====================
echo "### Updating compose profile"
declare -a MEDCO_NODES_URL OIDC_JWKS_URLS OIDC_JWT_ISSUERS OIDC_CLIENT_IDS OIDC_JWT_USER_ID_CLAIMS
declare -a WS_CLIENTS_LISTEN_ADDRESSES WS_CLIENTS_WS_SERVER_URLS WS_CLIENTS_WS_SERVER_PATH_PREFIXES WS_CLIENTS_DEST_ADDRESSES
pushd "${CONF_FOLDER}"

for ((i=0; i<NB_NODES; i++)); do
  CURR_NODE_IDX=$(printf "%03d" "$i")
  CURR_NODE_DNS_NAME=$(<"srv$CURR_NODE_IDX-nodednsname.txt")

  MEDCO_NODES_URL+=("https://$CURR_NODE_DNS_NAME/medco")
  OIDC_JWKS_URLS+=("https://$CURR_NODE_DNS_NAME/auth/realms/master/protocol/openid-connect/certs")
  OIDC_JWT_ISSUERS+=("https://$CURR_NODE_DNS_NAME/auth/realms/master")
  OIDC_CLIENT_IDS+=("medco")
  OIDC_JWT_USER_ID_CLAIMS+=("preferred_username")

  WS_CLIENTS_LISTEN_ADDRESSES+=("0.0.0.0:3$CURR_NODE_IDX")
  WS_CLIENTS_WS_SERVER_URLS+=("wss://$CURR_NODE_DNS_NAME:443")
  WS_CLIENTS_WS_SERVER_PATH_PREFIXES+=("unlynx")
  WS_CLIENTS_DEST_ADDRESSES+=("medco-unlynx:2001")
done
popd

export IFS=,
cat >> "${COMPOSE_FOLDER}/.env" <<EOF
MEDCO_NODES_URL=${MEDCO_NODES_URL[*]}
OIDC_JWKS_URLS=${OIDC_JWKS_URLS[*]}
OIDC_JWT_ISSUERS=${OIDC_JWT_ISSUERS[*]}
OIDC_CLIENT_IDS=${OIDC_CLIENT_IDS[*]}
OIDC_JWT_USER_ID_CLAIMS=${OIDC_JWT_USER_ID_CLAIMS[*]}

WS_CLIENTS_LISTEN_ADDRESSES=${WS_CLIENTS_LISTEN_ADDRESSES[*]}
WS_CLIENTS_WS_SERVER_URLS=${WS_CLIENTS_WS_SERVER_URLS[*]}
WS_CLIENTS_WS_SERVER_PATH_PREFIXES=${WS_CLIENTS_WS_SERVER_PATH_PREFIXES[*]}
WS_CLIENTS_DEST_ADDRESSES=${WS_CLIENTS_DEST_ADDRESSES[*]}
EOF
echo "### DONE! MedCo profile ready"
