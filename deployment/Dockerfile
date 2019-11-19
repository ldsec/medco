FROM golang:1.12 as build

COPY ./ /src

# compile and install medco-connector-server
WORKDIR /src
RUN CGO_ENABLED=0 go install -v ./cmd/medco-connector-server/...

# -------------------------------------------
FROM golang:1.12-alpine as release

COPY deployment/docker-entrypoint.sh /usr/local/bin/
RUN apk update && apk add bash && rm -rf /var/cache/apk/* && \
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
    JWKS_URL=http://keycloak:8080/auth/realms/master/protocol/openid-connect/certs \
    OIDC_JWT_ISSUER=http://keycloak:8080/auth/realms/master \
    OIDC_CLIENT_ID=medco \
    OIDC_JWT_USER_ID_CLAIM=preferred_username \
    MEDCO_OBFUSCATION_MIN=5 \
    MEDCO_NODES_URL="http://medco-connector-srv0/medco,http://medco-connector-srv1/medco,http://medco-connector-srv2/medco" \
    MEDCO_NODE_IDX=0 \
    MEDCO_OBFUSCATION_MIN=5 \
    I2B2_DB_HOST=postgresql \
    I2B2_DB_PORT=5432 \
    I2B2_DB_USER=postgres \
    I2B2_DB_PW=postgres \
    I2B2_DB_NAME=i2b2medco

EXPOSE 1999
ENTRYPOINT docker-entrypoint.sh medco-connector-server --write-timeout=${SERVER_HTTP_WRITE_TIMEOUT_SECONDS}s
