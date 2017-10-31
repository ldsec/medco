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
rm -f "$CONF_FOLDER"/*.keystore "$CONF_FOLDER"/shrine_downstream_nodes.conf "$CONF_FOLDER"/shrine_alias_map.conf "$CONF_FOLDER"/*.cer

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
    keytool -genkeypair -keysize 2048 -alias "$NODE_DNS-private" -validity 7300 \
        -dname "CN=$NODE_DNS, OU=$NODE_DNS, O=SHRINE Network, L=Lausanne, S=VD, C=CH" -ext "SAN=IP:$NODE_IP" \
        -keyalg RSA -keypass "$KEYSTORE_PW" -storepass "$KEYSTORE_PW" -keystore "$KEYSTORE"

    keytool -export -alias "$NODE_DNS-private" -storepass "$KEYSTORE_PW" -file "$CONF_FOLDER/$NODE_DNS.cer" -keystore "$KEYSTORE"

    # add entry in the downstream nodes and alias map
    echo "\"$NODE_DNS\" = \"https://$NODE_DNS:6443/shrine/rest/adapter/requests\"" >> "$CONF_FOLDER/shrine_downstream_nodes.conf"
    echo "\"$NODE_DNS\" = \"$NODE_DNS\"" >> "$CONF_FOLDER/shrine_alias_map.conf"

    #todo: unlynx keys

done

# import certificates of network nodes into the keystores
for KEYSTORE in "$CONF_FOLDER"/*.keystore
do
    for CERTIFICATE in "$CONF_FOLDER"/*.cer
    do
        NODE_DNS=$(basename "$CERTIFICATE" ".cer")
        #if [ "$NODE_DNS" != $(basename "$KEYSTORE" ".keystore") ]
        #then
            keytool -noprompt -import -v -trustcacerts -alias "$NODE_DNS" -file "$CERTIFICATE" -keystore "$KEYSTORE"  -keypass "$KEYSTORE_PW"  -storepass "$KEYSTORE_PW"
        #fi
    done

    keytool -list -v -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW"
done
