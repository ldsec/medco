#!/usr/bin/env bash

# ----------------------------------
# Global variables (editable)
# ----------------------------------

REQUIRED_PORTS="80,443,5432,2001,2002"


# ----------------------------------
# echo Colors
# ----------------------------------

NOCOLOR='\033[0m'
RED='\033[0;31m'
GREEN='\033[0;32m'


# ----------------------------------
# Functions
# ----------------------------------

die() {
    echo -e " ${RED}ERROR${NOCOLOR}"
    echo -e "${RED}$1${NOCOLOR}" >&2
    exit 1
}


# ----------------------------------
# Script
# ----------------------------------

for server_address in "$@"; do
    nmap_command=$(nmap --open $server_address | grep open)

    for port in $(echo $REQUIRED_PORTS | sed "s/,/ /g"); do
        if [ -z "$(echo "$nmap_command" | grep $port)" ]; then
            die "The port $port is not open on machine $server_address."
        fi
        echo "Port $port ok!"
    done
done

echo
echo -e "${GREEN}All MedCo ports have been successfully tested!${NOCOLOR}"

exit 0