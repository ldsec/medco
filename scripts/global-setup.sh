#!/usr/bin/env bash

# ----------------------------------
# Global variables (editable)
# ----------------------------------

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MEDCO_SETUP_DIR="~/medco_install_test"
MEDCO_SETUP_VER=v2.0.1
MEDCO_SETUP_NETWORK_NAME="medco-install-script-test"
MEDCO_REPO="https://github.com/ldsec/medco.git"
CLEANUP=1


# ----------------------------------
# Script variables (NOT editable)
# ----------------------------------

MEDCO_REPO_NAME=$(echo ${MEDCO_REPO%".git"} | awk -F/ '{print $NF}')
TOOLS_LIST="docker,docker-compose,git,nmap,openssl"
REQUIRED_PORTS="22" # TODO: Set ports 80,443,5432,2000,2001


# ----------------------------------
# echo Colors
# ----------------------------------

NOCOLOR='\033[0m'
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'


# ----------------------------------
# Formatting arguments to ssh command
# ----------------------------------

SSH_SERVER_LIST=()
for server_address in "$@"; do
    SSH_SERVER_LIST+=( "ssh $server_address" )
done


# ----------------------------------
# Functions
# ----------------------------------

print_step_header() {
    step_header_size=40
    number_of_char=$((step_header_size - (${#1} / 2)))
    echo -n -e $CYAN
    eval $(echo printf '=%.0s' {1..$((step_header_size * 2 + 3))})
    echo
    eval $(echo printf '=%.0s' {1..$number_of_char})
    echo -n " $1 "
    eval $(echo printf '=%.0s' {1..$number_of_char})
    echo
    eval $(echo printf '=%.0s' {1..$((step_header_size * 2 + 3))})
    echo
    echo -n -e $NOCOLOR
}

get_node_number() {
    if [ ${#1} -lt 3 ]; then
        number_of_char=$((3 - ${#1}))
        eval $(echo printf '0%.0s' {1..$((number_of_char))})
    fi
    echo $1
}

get_server_address() {
    echo $1 | awk '{print $NF}'
}

ok() {
    echo -e " ${GREEN}OK${NOCOLOR}"
}

die() {
    echo -e " ${RED}ERROR${NOCOLOR}"
    echo -e "${RED}$1${NOCOLOR}" >&2
    exit 1
}


# ----------------------------------
# Script
# ----------------------------------

print_step_header "Checking local SSH config (1/6)"
server_count=1
for server_ssh in "${SSH_SERVER_LIST[@]}"; do
    server_address=$(get_server_address "$server_ssh")
    echo "Checking SSH config of $server_address... ($server_count/$#)"

    echo -n "Test SSH Config..."
    if [ -z "$(cat ~/.ssh/config | grep -w "Host $server_address")" ]; then
        die "$server_address is not configured in ~/.ssh/config"
    fi
    ok
    
    echo -e "${GREEN}$(get_server_address "$server_ssh") is configured in SSH!${NOCOLOR}"
    ((server_count++))
done

print_step_header "Checking remote machines setup (2/6)"
server_count=1
for server_ssh in "${SSH_SERVER_LIST[@]}"; do
    echo "Testing $(get_server_address "$server_ssh")... ($server_count/$#)"

    echo -n "Initial connection..."
        initial_connection=$($server_ssh ls 2>/dev/null)
        if [ $? != 0 ]; then
            die "The machine $(get_server_address "$server_ssh") is not accessible."
        fi
    ok

    echo -n "Tools list..."
    for tool in $(echo $TOOLS_LIST | sed "s/,/ /g"); do
        if [ -z $($server_ssh which $tool) ]; then
            die "The machine $(get_server_address "$server_ssh") do not have $tool installed. Please install it manually first."
        fi
    done
    ok

    echo -n "Ports..."
    for server in "${SSH_SERVER_LIST[@]}"; do
        if [ "$server" != "$server_ssh" ]; then # Change != to = for testing
            server_address=$(get_server_address "$server")
            nmap_command=$($server_ssh "nmap --open $server_address | grep open")

            for port in $(echo $REQUIRED_PORTS | sed "s/,/ /g"); do
                if [ -z "$(echo "$nmap_command" | grep $port)" ]; then
                    die "The machine $(get_server_address "$server_ssh") cannot access $server_address on port $port."
                fi
            done
        fi
    done
    ok

    echo -e "${GREEN}$(get_server_address "$server_ssh") seems valid!${NOCOLOR}"
    ((server_count++))
done

print_step_header "Downloading MedCo (3/6)"
server_count=1
for server_ssh in "${SSH_SERVER_LIST[@]}"; do
    echo "Downloading MedCo on $(get_server_address "$server_ssh")... ($server_count/$#)"

    MEDCO_SETUP_NODE_IDX=$((server_count - 1))
    NODE_NUMBER=$(get_node_number $MEDCO_SETUP_NODE_IDX)

    if [ $CLEANUP = 1 ]; then
        echo -n "Cleanup MedCo repo..."
        $server_ssh "rm -rf $MEDCO_SETUP_DIR"
        ok
    fi

    echo -n "Cloning MedCo..."
    clone_command=$($server_ssh "git clone --branch $MEDCO_SETUP_VER $MEDCO_REPO $MEDCO_SETUP_DIR" 2>&1) # TODO: Add --depth 1 when PR#110 is merged and associated with a new tag
    $server_ssh "cd $MEDCO_SETUP_DIR && git checkout origin/add-custom-docker-network -- scripts/network-profile-tool/docker-compose.yml" 2>&1 # TODO: Remove when PR#110 is merged and associated with a new tag
    if [ $? != 0 ]; then
            die "The machine $(get_server_address "$server_ssh") cannot clone the MedCo repo.\n$clone_command"
    fi
    ok

    echo -e "${GREEN}MedCo downloaded on $(get_server_address "$server_ssh")!${NOCOLOR}"
    ((server_count++))
done

print_step_header "MedCo install step 1 (4/6)"
server_count=1
for server_ssh in "${SSH_SERVER_LIST[@]}"; do
    echo "MedCo install step 1 on $(get_server_address "$server_ssh")... ($server_count/$#)"

    echo -n "Launch step 1..."
    MEDCO_SETUP_NODE_IDX=$((server_count - 1))
    NODE_NUMBER=$(get_node_number $MEDCO_SETUP_NODE_IDX)
    step1_command=$($server_ssh "export MEDCO_SETUP_VER=$MEDCO_SETUP_VER && \
                    cd $MEDCO_SETUP_DIR/scripts/network-profile-tool && \
                    chmod +x step1.sh && \
                    echo y | ./step1.sh $MEDCO_SETUP_NETWORK_NAME $MEDCO_SETUP_NODE_IDX $(get_server_address "$server_ssh")" 2>&1)
    if [ $? != 0 ]; then
            die "Error on step 1 on $(get_server_address "$server_ssh").\n$step1_command"
    fi
    ok

    echo -n "Pulling configuration file..."
    rm -rf tmp_server_configuration
    mkdir tmp_server_configuration
    server_address=$(get_server_address "$server_ssh")
    configuration_path="$MEDCO_SETUP_DIR/deployments/network-$MEDCO_SETUP_NETWORK_NAME-node$NODE_NUMBER/configuration"
    scp_pull_command=$(scp $server_address:$configuration_path/srv$NODE_NUMBER-public.tar.gz tmp_server_configuration/ 2>&1)
    if [ $? != 0 ]; then
        die "Error on pulling configuration file on $(get_server_address "$server_ssh").\n$scp_pull_command"
    fi
    ok

    echo -n "Transfering file to other nodes..."
    for server in "${SSH_SERVER_LIST[@]}"; do
        if [ "$server" != "$server_ssh" ]; then # Change != to = for testing
            server_address=$(get_server_address "$server")
            $server "mkdir -p $configuration_path"
            scp_push_command=$(scp tmp_server_configuration/srv$NODE_NUMBER-public.tar.gz $server_address:$configuration_path/srv$NODE_NUMBER-public.tar.gz)
            if [ $? != 0 ]; then
                die "Error on pushing configuration file to $(get_server_address "$server").\n$scp_push_command"
            fi
        fi
    done
    ok

    echo -e "${GREEN}Step 1 finished and configuration file of $(get_server_address "$server_ssh") transfered!${NOCOLOR}"
    ((server_count++))
done

rm -rf tmp_server_configuration

print_step_header "MedCo install step 2 (5/6)"
server_count=1
for server_ssh in "${SSH_SERVER_LIST[@]}"; do
    echo "MedCo install step 2 on $(get_server_address "$server_ssh")... ($server_count/$#)"

    echo -n "Launch step 2..."
    MEDCO_SETUP_NODE_IDX=$((server_count - 1))
    NODE_NUMBER=$(get_node_number $MEDCO_SETUP_NODE_IDX)
    step2_command=$($server_ssh "export MEDCO_SETUP_VER=$MEDCO_SETUP_VER && \
                    cd $MEDCO_SETUP_DIR/scripts/network-profile-tool && \
                    chmod +x step2.sh && \
                    echo y | ./step2.sh $MEDCO_SETUP_NETWORK_NAME $MEDCO_SETUP_NODE_IDX" 2>&1)
    if [ $? != 0 ]; then
            die "Error on step 2 on $(get_server_address "$server_ssh").\n$step2_command"
    fi
    ok

    echo -e "${GREEN}Step 2 finished for $(get_server_address "$server_ssh")!${NOCOLOR}"
    ((server_count++))
done

print_step_header "Deployment (6/6)"
server_count=1
for server_ssh in "${SSH_SERVER_LIST[@]}"; do
    echo "Deployment on $(get_server_address "$server_ssh")... ($server_count/$#)"

    echo -n "Launch deployment..."
    MEDCO_SETUP_NODE_IDX=$((server_count - 1))
    NODE_NUMBER=$(get_node_number $MEDCO_SETUP_NODE_IDX)
    deployment_command=$($server_ssh "cd $MEDCO_SETUP_DIR/deployments/network-$MEDCO_SETUP_NETWORK_NAME-node$NODE_NUMBER && \
                        make down && \
                        make pull && \
                        make up" 2>&1)
    if [ $? != 0 ]; then
            die "Error on deployment on $(get_server_address "$server_ssh").\n$deployment_command"
    fi
    ok

    echo -e "${GREEN}Deployment finished for $(get_server_address "$server_ssh")!${NOCOLOR}"
    ((server_count++))
done

echo
echo -e "${GREEN}All MedCo servers have been successfully installed and deployed!${NOCOLOR}"

exit 0