#!/usr/bin/env bash
set -Eeo pipefail

# trust the certificates of other nodes
if [[ -d /medco-configuration ]]; then
  NB_CA_CERTS=$(find /medco-configuration -maxdepth 1 -name '*.crt' | wc -l)
  if [[ "$NB_CA_CERTS" != 0 ]]; then
    cp -f /medco-configuration/*.crt /usr/local/share/ca-certificates/
    update-ca-certificates
    echo "WS TUNNEL: $NB_CA_CERTS CA certificates added"
  else
    echo "WS TUNNEL: no CA certificate added"
  fi
fi

# run server
if [[ -n "$SERVER_LISTEN_URL" ]]; then
  /wstunnel --server "$SERVER_LISTEN_URL" --restrictTo="$SERVER_DEST_ADDRESS" &
  echo "WS TUNNEL: ran server on $SERVER_LISTEN_URL"
else
  echo "WS TUNNEL: no server was ran"
fi

# run clients
if [[ -n "$CLIENTS_LISTEN_ADDRESSES" ]]; then
  IFS=',' read -ra LISTEN_ADDRESSES_ARR <<< "$CLIENTS_LISTEN_ADDRESSES"
  IFS=',' read -ra WS_SERVER_URLS_ARR <<< "$CLIENTS_WS_SERVER_URLS"
  IFS=',' read -ra WS_SERVER_PATH_PREFIXES_ARR <<< "$CLIENTS_WS_SERVER_PATH_PREFIXES"
  IFS=',' read -ra DEST_ADDRESSES_ARR <<< "$CLIENTS_DEST_ADDRESSES"

  NB_CLIENTS=${#LISTEN_ADDRESSES_ARR[@]}
  echo "WS TUNNEL: about to start $NB_CLIENTS clients..."
  for ((i=0; i<NB_CLIENTS; i++)); do
    LISTEN_ADDRESS=${LISTEN_ADDRESSES_ARR[i]}
    WS_SERVER_URL=${WS_SERVER_URLS_ARR[i]}
    WS_SERVER_PATH_PREFIX=${WS_SERVER_PATH_PREFIXES_ARR[i]}
    DEST_ADDRESS=${DEST_ADDRESSES_ARR[i]}

    /wstunnel -L "$LISTEN_ADDRESS:$DEST_ADDRESS" \
      --upgradePathPrefix="$WS_SERVER_PATH_PREFIX" \
      "$WS_SERVER_URL" &
    echo "WS TUNNEL: ran client on $LISTEN_ADDRESS"
  done
else
  echo "WS TUNNEL: no client was ran"
fi

# wait on subshells and return exit code
wait -n
exit $?
