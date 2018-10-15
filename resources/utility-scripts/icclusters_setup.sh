#!/usr/bin/env bash

SSH_TYPE="-t ssh-ed25519"
SERVERS="$@"
SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

for s in $SERVERS; do
    echo "Setting up iccluster$s.iccluster.epfl.ch..."
    login=root@iccluster$s.iccluster.epfl.ch
    #cat install_script.sh | ssh $login /bin/bash
    cat "$SCRIPT_FOLDER"/ubuntu_prereqs_setup.sh | sshpass -p "1" ssh -o StrictHostKeyChecking=no $login /bin/bash
done
