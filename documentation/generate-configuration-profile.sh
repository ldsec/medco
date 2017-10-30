#!/bin/bash
set -e
shopt -s nullglob

# usage: bash generate-configuration-profile.sh CONFIGURATION_FOLDER KEYSTORE_PASSWORD NODE_DNS_1 NODE_IP_1 NODE_DNS_2 NODE_IP_3 NODE_DNS_3 NODE_IP_3 ...
if [ $# -lt 4 ]
then
    echo "Wrong number of arguments, usage: bash generate-configuration-profile.sh CONFIGURATION_FOLDER KEYSTORE_PASSWORD NODE_DNS_1 NODE_IP_1 NODE_DNS_2 NODE_IP_3 NODE_DNS_3 NODE_IP_3 ..."
    exit
fi

# arguments
CONF_FOLDER="$1"
KEYSTORE_PW="$2"
shift
shift

# clean up previous entries
mkdir -p "$CONF_FOLDER"
rm -f "$CONF_FOLDER"/*.keystore "$CONF_FOLDER"/*.conf "$CONF_FOLDER"/*.cer

# generate private and keystore for each node
NODE_IDX="-1"
while [ $# -gt 0 ]
do
    NODE_DNS="$1"
    NODE_IP="$2"
    NODE_IDX=$((NODE_IDX+1))
    KEYSTORE="$CONF_FOLDER/$NODE_DNS.keystore"
    shift
    shift

    # generate the node certificate in the keystore and export it
    keytool -genkeypair -keysize 2048 -alias "$NODE_DNS" -validity 7300 \
        -dname "CN=$NODE_DNS, OU=$NODE_DNS, O=SHRINE Network, L=Lausanne, S=VD, C=CH" -ext "SAN=IP:$NODE_IP" \
        -keyalg RSA -keypass "$KEYSTORE_PW" -storepass "$KEYSTORE_PW" -keystore "$KEYSTORE"

    keytool -export -alias "$NODE_DNS" -storepass "$KEYSTORE_PW" -file "$CONF_FOLDER/$NODE_DNS.cer" -keystore "$KEYSTORE"

    # add entry in the downstream nodes and alias map
    echo "\"$NODE_DNS\" = \"$NODE_DNS\"" >> "$CONF_FOLDER/shrine_aliasMap.conf"

    #todo: unlynx keys

done

# import certificates of network nodes into the keystores
for KEYSTORE in "$CONF_FOLDER"/*.keystore
do
    for CERTIFICATE in "$CONF_FOLDER"/*.cer
    do
        OTHER_NODE_DNS=$(basename "$CERTIFICATE" ".cer")
        CURRENT_NODE_DNS=$(basename "$KEYSTORE" ".keystore")
        if [ "$OTHER_NODE_DNS" != "$CURRENT_NODE_DNS" ]
        then
            keytool -noprompt -import -v -trustcacerts -alias "$OTHER_NODE_DNS" -file "$CERTIFICATE" -keystore "$KEYSTORE"  -keypass "$KEYSTORE_PW"  -storepass "$KEYSTORE_PW"

            # generate aliasMap and downstreamNodes
            echo "\"$OTHER_NODE_DNS\" = \"https://$OTHER_NODE_DNS:6443/shrine/rest/adapter/requests\"" >> "$CONF_FOLDER/${CURRENT_NODE_DNS}_downstreamNodes.conf"
        fi
    done

    keytool -list -v -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW"
done
