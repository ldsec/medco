FROM golang:1.12.10 as build

COPY ./ /src

# compile and install medco-cli-client
WORKDIR /src
RUN CGO_ENABLED=0 go install -v ./cmd/medco-cli-client/...

# -------------------------------------------
FROM golang:1.12.5-alpine as release

COPY deployment/docker-entrypoint.sh /usr/local/bin/
RUN apk update && apk add bash && rm -rf /var/cache/apk/* && \
    chmod a+x /usr/local/bin/docker-entrypoint.sh

COPY --from=build /go/bin/medco-cli-client /go/bin/

# run-time environment
ENV LOG_LEVEL=5 \
    UNLYNX_GROUP_FILE_PATH=/medco-configuration/group.toml \
    MEDCO_NODE_IDX=0 \
    CLIENT_QUERY_TIMEOUT_SECONDS=660 \
    CLIENT_GENOMIC_ANNOTATIONS_QUERY_TIMEOUT_SECONDS=10 \
    MEDCO_CONNECTOR_URL=http://medco-connector-srv0/medco \
    OIDC_REQ_TOKEN_URL=http://keycloak:8080/auth/realms/master/protocol/openid-connect/token \
    OIDC_REQ_TOKEN_CLIENT_ID=medco

ENTRYPOINT ["docker-entrypoint.sh", "medco-cli-client"]
