FROM golang:1.11 as build

COPY ./ /src

# compile and install medco-connector
WORKDIR /src
RUN CGO_ENABLED=0 go build -v ./... && \
    CGO_ENABLED=0 go install -v ./...

# -------------------------------------------
FROM golang:1.11-alpine as release

COPY --from=build /go/bin/medco-connector-server /go/bin/

# swagger server configuration
ENV HOST=0.0.0.0 \
    PORT=1999

# run-time environment
ENV PICSURE2_INTERNAL_TOKEN=xxx \
    I2B2_HIVE_URL=http://i2b2:8080/i2b2/services \
    I2B2_LOGIN_DOMAIN=i2b2medco \
    I2B2_LOGIN_PROJECT=MedCo \
    I2B2_LOGIN_USER=e2etest \
    I2B2_LOGIN_PASSWORD=e2etest \
    LOG_LEVEL=5

EXPOSE 1999
ENTRYPOINT ["medco-connector-server"]
