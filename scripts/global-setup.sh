#!/usr/bin/env bash

# ----------------------------------
# Global variables
# ----------------------------------

SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
MEDCO_INSTALL_DIR="~/script_tests"
TOOLS_LIST="docker,docker-compose,git,nmap"
REQUIRED_PORTS="22" # TODO: Set ports 80,443,5432,2000,2001
MEDCO_REPO="https://github.com/ldsec/medco.git"


# ----------------------------------
# echo Colors
# ----------------------------------

NOCOLOR='\033[0m'
RED='\033[0;31m'
GREEN='\033[0;32m'


print_step_header() {
    step_header_size=40
    number_of_char=$((step_header_size - (${#1} / 2)))
    eval $(echo printf '=%.0s' {1..$number_of_char})
    echo -n " $1 "
    eval $(echo printf '=%.0s' {1..$number_of_char})
    echo
}

get_server_address() {
    echo $1 | awk '{print $NF}'
}

die() {
    echo -e " ${RED}ERROR${NOCOLOR}"
    echo -e "${RED}$1${NOCOLOR}" >&2
    exit 1
}

print_step_header "Checking machines setup (1/3)"
server_count=1
for server_ssh in "$@"; do
    echo "Testing $(get_server_address "$server_ssh")... ($server_count/$#)"

    echo -n "Testing tools list..."
    for tool in $(echo $TOOLS_LIST | sed "s/,/ /g"); do
        if [ -z $($server_ssh which $tool) ]; then
            die "The machine at $server_ssh do not have $tool installed. Please install it manually first."
        fi
    done
    echo -e " ${GREEN}OK${NOCOLOR}"

    echo -n "Testing ports..."
    for server in "$@"; do
        if [ "$server" = "$server_ssh" ]; then # TODO: change to !=
            remote_server_address=$(get_server_address "$server")
            nmap_result=$($server_ssh "nmap --open $remote_server_address | grep open")

            for port in $(echo $REQUIRED_PORTS | sed "s/,/ /g"); do
                if [ -z "$(echo "$nmap_result" | grep $port)" ]; then
                    die "The machine at $server_ssh cannot access $remote_server_address on port $port."
                fi
            done
        fi
    done
    echo -e " ${GREEN}OK${NOCOLOR}"

    echo -e "${GREEN}$(get_server_address "$server_ssh") seems valid!${NOCOLOR}"
    ((server_count++))
done

print_step_header "Downloading MedCo (2/3)"
server_count=1
for server_ssh in "$@"; do
    echo "Downloading MedCo on $(get_server_address "$server_ssh")... ($server_count/$#)"

    $server_ssh "mkdir -p $MEDCO_INSTALL_DIR && \
                cd $MEDCO_INSTALL_DIR && \
                rm -rf medco && \
                git clone $MEDCO_REPO && \
                cd medco && \
                ls"

    # launch steps1/2/3 one by one

    die finished

    echo -e "${GREEN}MedCo downloaded!${NOCOLOR}"
    ((server_count++))
done



