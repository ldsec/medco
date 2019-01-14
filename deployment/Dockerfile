FROM medco/unlynx:v1.2a-alpha as build

# build time environment variables
ENV MEDCO_UNLYNX_REPO="github.com/lca1/medco-unlynx"

# adjust timeout values
RUN sed -i 's/conn.SetReadDeadline(time.Now().Add(5 \* time.Minute))/conn.SetReadDeadline(time.Now().Add(5 \* time.Hour))/' /go/src/github.com/dedis/onet/websocket.go && \
    sed -i 's/const expirationTime = 1 \* time.Minute/const expirationTime = 1 \* time.Hour/' /go/src/github.com/dedis/onet/overlay.go

COPY ./ /go/src/$MEDCO_UNLYNX_REPO

# get remaining dependencies, compile and install medco-unlynx binary
# CGO_ENABLED=0 in order to be able to run from alpine
WORKDIR /go/src/$MEDCO_UNLYNX_REPO
RUN go get -v -d ./... && \
    CGO_ENABLED=0 go build -v ./... && \
    CGO_ENABLED=0 go install -v ./... && \
    mv /go/bin/app /go/bin/medco-unlynx

# -------------------------------------------
FROM golang:1.11-alpine as release

# run time environment variables
ENV NODE_IDX="0" \
    UNLYNX_DEBUG_LEVEL="1" \
    MEDCO_CONF_DIR="/medco-configuration"
ENV UNLYNX_BIN_EXPORT_PATH="$MEDCO_CONF_DIR/medco-unlynx"

COPY --from=build /go/bin/medco-unlynx /go/bin/
COPY deployment/docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh && \
    apk add --no-cache bash outils-md5

VOLUME "$MEDCO_CONF_DIR"
ENTRYPOINT docker-entrypoint.sh
