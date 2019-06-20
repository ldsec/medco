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
    PORT=1999

# run-time environment
ENV I2B2_HIVE_URL=http://i2b2:8080/i2b2/services \
    I2B2_LOGIN_DOMAIN=i2b2medco \
    I2B2_LOGIN_PROJECT=MedCo \
    I2B2_LOGIN_USER=e2etest \
    I2B2_LOGIN_PASSWORD=e2etest \
    LOG_LEVEL=5 \
    UNLYNX_GROUP_FILE_PATH=/medco-configuration/group.toml \
    UNLYNX_GROUP_FILE_IDX=0 \
    JWKS_URL=http://keycloak:8080/auth/realms/master/protocol/openid-connect/certs \
    OIDC_JWT_ISSUER=http://keycloak:8080/auth/realms/master \
    OIDC_CLIENT_ID=medco \
    OIDC_JWT_USER_ID_CLAIM=preferred_username \
    MEDCO_OBFUSCATION_MIN_VARIANCE=5

EXPOSE 1999
ENTRYPOINT ["docker-entrypoint.sh", "medco-connector-server"]
