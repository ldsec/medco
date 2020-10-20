FROM golang:1.14 as build
ARG MEDCO_VERSION=dev

COPY ./ /src
WORKDIR /src

# get dependencies
RUN go get -v -d ./...

# hack to adjust timeout values in onet
RUN cd /go/pkg/mod/go.dedis.ch/onet && \
    chmod u+w -R . && \
    find . -maxdepth 2 -name overlay.go | xargs sed -i \
        's/const expirationTime = 1 \* time.Minute/const expirationTime = 1 \* time.Hour/' && \
    find . -maxdepth 2 -name websocket.go | xargs sed -i \
        's/conn.SetReadDeadline(time.Now().Add(5 \* time.Minute))/conn.SetReadDeadline(time.Now().Add(5 \* time.Hour))/'

# compile and install all binaries
RUN CGO_ENABLED=0 go install -v -ldflags "-X github.com/ldsec/medco.Version=${MEDCO_VERSION}" ./...

# -------------------------------------------
FROM golang:1.14-alpine as release

COPY build/package/medco/docker-entrypoint.sh /usr/local/bin/
RUN apk update && apk add bash postgresql && rm -rf /var/cache/apk/* && \
    chmod a+x /usr/local/bin/docker-entrypoint.sh
COPY --from=build /go/bin /go/bin

# medco-connector-server, medco-cli-client, medco-loader
ENV LOG_LEVEL=5 \
    UNLYNX_GROUP_FILE_PATH=/medco-configuration/group.toml \
    MEDCO_NODE_IDX=0

# medco-connector, medco-loader
ENV I2B2_DB_HOST=postgresql \
    I2B2_DB_PORT=5432 \
    I2B2_DB_NAME=i2b2medco \
    I2B2_DB_USER=i2b2 \
    I2B2_DB_PASSWORD=i2b2 \
    MC_DB_HOST=postgresql \
    MC_DB_PORT=5432 \
    MC_DB_NAME=medcoconnector \
    MC_DB_USER=medcoconnector \
    MC_DB_PASSWORD=medcoconnector

# medco-connector-server
ENV HOST=0.0.0.0 \
    PORT=1999 \
    SERVER_HTTP_WRITE_TIMEOUT_SECONDS=600 \
    I2B2_HIVE_URL=http://i2b2:8080/i2b2/services \
    I2B2_LOGIN_DOMAIN=i2b2medco \
    I2B2_LOGIN_PROJECT=MedCo \
    I2B2_LOGIN_USER=e2etest \
    I2B2_LOGIN_PASSWORD=e2etest \
    I2B2_WAIT_TIME_SECONDS=300 \
    UNLYNX_TIMEOUT_SECONDS=300 \
    OIDC_JWKS_URLS=http://keycloak:8080/auth/realms/master/protocol/openid-connect/certs \
    OIDC_JWT_ISSUERS=http://keycloak:8080/auth/realms/master \
    OIDC_CLIENT_IDS=medco \
    OIDC_JWT_USER_ID_CLAIMS=preferred_username \
    OIDC_JWKS_TIMEOUT_SECONDS=30 \
    MEDCO_OBFUSCATION_MIN=5 \
    MEDCO_NODES_URL="http://medco-connector-srv0/medco,http://medco-connector-srv1/medco,http://medco-connector-srv2/medco" \
    MEDCO_OBFUSCATION_MIN=5

# medco-unlynx
ENV UNLYNX_DEBUG_LEVEL="1"

# medco-cli-client
ENV CLIENT_SEARCH_TIMEOUT_SECONDS=10 \
    CLIENT_QUERY_TIMEOUT_SECONDS=660 \
    CLIENT_GENOMIC_ANNOTATIONS_QUERY_TIMEOUT_SECONDS=10 \
    MEDCO_CONNECTOR_URL=http://medco-connector-srv0/medco \
    OIDC_REQ_TOKEN_URL=http://keycloak:8080/auth/realms/master/protocol/openid-connect/token \
    OIDC_REQ_TOKEN_CLIENT_ID=medco

VOLUME "/medco-configuration"
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["ls -1 /go/bin"]
