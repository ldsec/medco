#!/bin/bash
set -e
shopt -s nullglob

# usage: bash generate-configuration-profile.sh CONFIGURATION_FOLDER KEYSTORE_PASSWORD NODE_DNS_1 NODE_IP_1 NODE_DNS_2 NODE_IP_2 NODE_DNS_3 NODE_IP_3 ...
if [ $# -lt 5 ]
then
    echo "Wrong number of arguments, usage: bash generate-configuration-profile.sh CONFIGURATION_FOLDER KEYSTORE_PASSWORD NODE_DNS_1 NODE_IP_1 NODE_DNS_2 NODE_IP_2 NODE_DNS_3 NODE_IP_3 ..."
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
rm -f "$CONF_FOLDER"/*.keystore "$CONF_FOLDER"/shrine_downstream_nodes.conf "$CONF_FOLDER"/*.cer

# generate private and keystore for each node
NODE_IDX="-1"
while [ $# -gt 0 ]
do
    NODE_DNS="$1"
    NODE_IP="$2"
    shift
    shift

    NODE_IDX=$((NODE_IDX+1))
    KEYSTORE="$CONF_FOLDER/srv$NODE_IDX.keystore"
    KEYSTORE_PRIVATE_ALIAS="srv$NODE_IDX-private"

    # generate node java keystore pair of keys
    keytool -genkeypair -keysize 2048 -alias "$KEYSTORE_PRIVATE_ALIAS" -validity 7300 \
        -dname "CN=$NODE_DNS, OU=LCA1, O=EPFL, L=Lausanne, S=VD, C=CH" \
        -ext "SAN=DNS:$NODE_DNS,IP:$NODE_IP" \
        -keyalg RSA -keypass "$KEYSTORE_PW" -storepass "$KEYSTORE_PW" -keystore "$KEYSTORE"

    # generate certificate signature request and sigh it with CA
    keytool -certreq -alias "$KEYSTORE_PRIVATE_ALIAS" -keyalg RSA -file "$SCRIPT_FOLDER/newreq.pem" -keypass "$KEYSTORE_PW" \
        -storepass "$KEYSTORE_PW" -keystore "$KEYSTORE" -ext "SAN=DNS:$NODE_DNS,IP:$NODE_IP"
    cat > "$SCRIPT_FOLDER/openssl.ext.tmp.cnf" <<EOL
        basicConstraints=CA:FALSE
        subjectAltName=@alt_names
        subjectKeyIdentifier = hash

        [ alt_names ]
        IP.1 = $NODE_IP
        DNS.1 = $NODE_DNS
EOL

    SSLEAY_CONFIG="-extfile $SCRIPT_FOLDER/openssl.ext.tmp.cnf" "$SCRIPT_FOLDER"/CA.sh -sign
    rm "$SCRIPT_FOLDER/newreq.pem" "$SCRIPT_FOLDER/openssl.ext.tmp.cnf"

    # import CA certificate and own certificate signed by CA (chained to the private key)
    keytool -noprompt -import -v -alias shrine-hub-ca -file "$SCRIPT_FOLDER"/CA/cacert.pem -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW"
    keytool -noprompt -import -v -alias "$KEYSTORE_PRIVATE_ALIAS" -file "$SCRIPT_FOLDER"/newcert.pem -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW" \
        -keypass "$KEYSTORE_PW" -trustcacerts

    # lighttpd certificates
    keytool -noprompt -importkeystore -srckeystore "$KEYSTORE" -srcalias "$KEYSTORE_PRIVATE_ALIAS" -destkeystore "$KEYSTORE".p12 \
        -deststoretype PKCS12 -srcstorepass "$KEYSTORE_PW" -deststorepass "$KEYSTORE_PW"
    openssl pkcs12 -in "$KEYSTORE".p12 -out "$CONF_FOLDER/srv$NODE_IDX.pem" -password pass:"$KEYSTORE_PW" -nodes
    #cat "$SCRIPT_FOLDER"/newcert.pem "$CONF_FOLDER/srv$NODE_IDX-private.pem" > "$CONF_FOLDER/srv$NODE_IDX.pem"
    rm "$KEYSTORE".p12 #"$CONF_1FOLDER/srv$NODE_IDX-private.pem"

    #todo: remove inermediate files

    #keytool -import -v -alias shrine-hub-https -file "$SCRIPT_FOLDER" -keystore $KEYSTORE_FILE -storepass $KEYSTORE_PASSWORD
    #keytool -export -alias "$NODE_DNS-private" -storepass "$KEYSTORE_PW" -file "$CONF_FOLDER/$NODE_DNS.cer" -keystore "$KEYSTORE"

    # add entry in the downstream nodes and alias map
    echo "\"Hospital $NODE_IDX\" = \"https://$NODE_DNS:6443/shrine/rest/adapter/requests\"" >> "$CONF_FOLDER/shrine_downstream_nodes.conf"

    #todo: unlynx keys
    #todo: cleanup

    keytool -list -v -keystore "$KEYSTORE" -storepass "$KEYSTORE_PW"
done

#CA generation
#CATOP=/home/misbach/repositories/medco-deployment/configuration-profiles/dev-3nodes-samehost/CA/ /etc/pki/tls/misc/CA -newca
#cartificate name: dev-3nodes-samehost
# cp cacert.pem in conf folder // todo: delete at beginning + delete toml and do unlynx keys
