FROM golang:1.12.5 as build

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
    UNLYNX_GROUP_FILE_IDX=0 \
    OIDC_CLIENT_ID=medco \
    CLIENT_QUERY_TIMEOUT_SECONDS=660 \
    PICSURE2_API_HOST=picsure:8080/pic-sure-api-2/PICSURE \
    PICSURE2_API_BASE_PATH="" \
    PICSURE2_API_SCHEME=https \
    PICSURE2_RESOURCES=MEDCO_testnetwork_0_a,MEDCO_testnetwork_1_b,MEDCO_testnetwork_2_c \
    OIDC_REQ_TOKEN_URL=http://keycloak:8080/auth/realms/master/protocol/openid-connect/token

ENTRYPOINT ["docker-entrypoint.sh", "medco-cli-client"]
