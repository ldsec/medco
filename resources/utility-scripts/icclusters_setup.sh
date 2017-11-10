#!/usr/bin/env bash

SSH_TYPE="-t ssh-ed25519"
SERVERS="$@"

for s in $SERVERS; do
    echo "Setting up iccluster$s.iccluster.epfl.ch..."
    login=root@iccluster$s.iccluster.epfl.ch
    #cat install_script.sh | ssh $login /bin/bash
    cat install_script.sh | sshpass -p "1" ssh -o StrictHostKeyChecking=no $login /bin/bash
done


