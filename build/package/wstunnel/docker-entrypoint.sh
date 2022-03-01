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
for i in {0..999}; do
  NODE_IDX=$(printf "%03d" "$i")

  LISTEN_ADDRESS_VAR=CLIENT_${NODE_IDX}_LISTEN_ADDRESS
  WS_SERVER_URL_VAR=CLIENT_${NODE_IDX}_WS_SERVER_URL
  WS_SERVER_PATH_PREFIX_VAR=CLIENT_${NODE_IDX}_WS_SERVER_PATH_PREFIX
  DEST_ADDRESS_VAR=CLIENT_${NODE_IDX}_DEST_ADDRESS

  if [[ -n "${!LISTEN_ADDRESS_VAR}" ]]; then
    /wstunnel -L "${!LISTEN_ADDRESS_VAR}:${!DEST_ADDRESS_VAR}" \
      --upgradePathPrefix="${!WS_SERVER_PATH_PREFIX_VAR}" \
      "${!WS_SERVER_URL_VAR}" &
    echo "WS TUNNEL: ran client on ${!LISTEN_ADDRESS_VAR}"
  fi
done

# wait on subshells and return exit code
wait -n
exit $?
