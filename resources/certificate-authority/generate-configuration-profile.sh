#!/bin/bash
set -e
shopt -s nullglob

# usage: bash generate-configuration-profile.sh CONFIGURATION_FOLDER KEYSTORE_PASSWORD NODE_DNS_1 NODE_IP_1 NODE_DNS_2 NODE_IP_3 NODE_DNS_3 NODE_IP_3 ...
if [ $# -lt 4 ]
then
    echo "Wrong number of arguments, usage: bash generate-configuration-profile.sh CONFIGURATION_FOLDER KEYSTORE_PASSWORD NODE_DNS_1 NODE_IP_1 NODE_DNS_2 NODE_IP_3 NODE_DNS_3 NODE_IP_3 ..."
    exit
fi

# variables
SCRIPT_FOLDER="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

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

    # generate node pair of keys
    keytool -genkeypair -keysize 2048 -alias "$NODE_DNS-private" -validity 7300 \
        -dname "CN=$NODE_DNS, OU=LCA1, O=EPFL, L=Lausanne, S=VD, C=CH" \
        -keyalg RSA -keypass "$KEYSTORE_PW" -storepass "$KEYSTORE_PW" -keystore "$KEYSTORE"

    # generate certificate signature request
    keytool -certreq -alias "$NODE_DNS-private" -keyalg RSA -file "$SCRIPT_FOLDER/newreq.pem" -keypass "$KEYSTORE_PW" \
        -storepass "$KEYSTORE_PW" -keystore "$KEYSTORE"

    # sign with the CA
    #mv "$CONF_FOLDER/$NODE_DNS.csr" "$CONF_FOLDER/newreq.pem"
    "$SCRIPT_FOLDER"/CA.sh -sign

    # import certificate of CA
    keytool -import -v -alias shrine-hub-ca -file "$SCRIPT_FOLDER"/CA/cacert.pem -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW"

    # import own certificate signed by CA to the private key
    keytool -import -v -alias "$NODE_DNS-private" -file "$SCRIPT_FOLDER"/newcert.pem -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW" \
        -keypass "$KEYSTORE_PW" -trustcacerts

    #keytool -import -v -alias shrine-hub-https -file "$SCRIPT_FOLDER" -keystore $KEYSTORE_FILE -storepass $KEYSTORE_PASSWORD
    #keytool -export -alias "$NODE_DNS-private" -storepass "$KEYSTORE_PW" -file "$CONF_FOLDER/$NODE_DNS.cer" -keystore "$KEYSTORE"

    # add entry in the downstream nodes and alias map
    echo "\"$NODE_DNS\" = \"https://$NODE_DNS:6443/shrine/rest/adapter/requests\"" >> "$CONF_FOLDER/shrine_downstream_nodes.conf"
    echo "\"$NODE_DNS\" = \"$NODE_DNS\"" >> "$CONF_FOLDER/shrine_alias_map.conf"

    #todo: unlynx keys

done

# import certificates of network nodes into the keystores
#for KEYSTORE in "$CONF_FOLDER"/*.keystore
#do
#    for CERTIFICATE in "$CONF_FOLDER"/*.cer
#    do
#        NODE_DNS=$(basename "$CERTIFICATE" ".cer")
        #if [ "$NODE_DNS" != $(basename "$KEYSTORE" ".keystore") ]
        #then
#            keytool -noprompt -import -v -trustcacerts -alias "$NODE_DNS" -file "$CERTIFICATE" -keystore "$KEYSTORE"  -keypass "$KEYSTORE_PW"  -storepass "$KEYSTORE_PW"
        #fi
#    done

#    keytool -list -v -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW"
#done

#CA generation
#CATOP=/home/misbach/repositories/medco-deployment/configuration-profiles/dev-3nodes-samehost/CA/ /etc/pki/tls/misc/CA -newca
#cartificate name: dev-3nodes-samehost