FROM golang:1.13 as build

COPY ./ /src

# compile and install medco-connector-server
WORKDIR /src
RUN CGO_ENABLED=0 go install -v ./cmd/medco-connector-server/...

# -------------------------------------------
FROM golang:1.13-alpine as release

COPY deployment/docker-entrypoint.sh /usr/local/bin/
RUN mkdir /usr/local/bin/querytools
COPY deployment/querytools /usr/local/bin/querytools
RUN apk update && apk add bash postgresql && rm -rf /var/cache/apk/* && \
    chmod a+x /usr/local/bin/docker-entrypoint.sh

COPY --from=build /go/bin/medco-connector-server /go/bin/

# swagger server configuration
ENV HOST=0.0.0.0 \
    PORT=1999 \
    SERVER_HTTP_WRITE_TIMEOUT_SECONDS=600

# run-time environment
ENV I2B2_HIVE_URL=http://i2b2:8080/i2b2/services \
    I2B2_LOGIN_DOMAIN=i2b2medco \
    I2B2_LOGIN_PROJECT=MedCo \
    I2B2_LOGIN_USER=e2etest \
    I2B2_LOGIN_PASSWORD=e2etest \
    I2B2_WAIT_TIME_SECONDS=300 \
    LOG_LEVEL=5 \
    UNLYNX_GROUP_FILE_PATH=/medco-configuration/group.toml \
    UNLYNX_TIMEOUT_SECONDS=300 \
    OIDC_JWKS_URLS=http://keycloak:8080/auth/realms/master/protocol/openid-connect/certs \
    OIDC_JWT_ISSUERS=http://keycloak:8080/auth/realms/master \
    OIDC_CLIENT_IDS=medco \
    OIDC_JWT_USER_ID_CLAIMS=preferred_username \
    OIDC_JWKS_TIMEOUT_SECONDS=30 \
    MEDCO_OBFUSCATION_MIN=5 \
    MEDCO_NODES_URL="http://medco-connector-srv0/medco,http://medco-connector-srv1/medco,http://medco-connector-srv2/medco" \
    MEDCO_NODE_IDX=0 \
    MEDCO_OBFUSCATION_MIN=5 \
    MC_DB_HOST=postgresql \
    MC_DB_PORT=5432 \
    MC_DB_USER=medcoconnector \
    MC_DB_PW=medcoconnector \
    MC_DB_NAME=medcoconnector \
    I2B2_DB_HOST=postgresql \
    I2B2_DB_PORT=5432 \
    I2B2_DB_USER=i2b2 \
    I2B2_DB_PW=i2b2 \
    I2B2_DB_NAME=i2b2medco

EXPOSE 1999
ENTRYPOINT docker-entrypoint.sh medco-connector-server --write-timeout=${SERVER_HTTP_WRITE_TIMEOUT_SECONDS}s
