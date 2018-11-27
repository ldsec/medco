FROM golang:1.11 as dev

# build time environment variables
ENV UNLYNX_VERSION="tags/v1.2-alpha" \
    MEDCO_UNLYNX_REPO="github.com/lca1/medco-unlynx" \
    MEDCO_UNLYNX_VERSION="v0.1"

# get maximum of dependencies and cache them
RUN go get -v -d github.com/dedis/onet/... && \
    go get -v -d github.com/dedis/kyber/... && \
    go get -v -d gopkg.in/urfave/cli.v1/... && \
    go get -v -d github.com/Knetic/govaluate/... && \
    go get -v -d github.com/btcsuite/goleveldb/... && \
    go get -v -d github.com/r0fls/gostats/... && \
    go get -v -d github.com/fanliao/go-concurrentMap/... && \
    go get -v -d github.com/lib/pq/... && \
    go get -v -d github.com/lca1/unlynx/...

# adjust timeout values and checkout specific version of unlynx
RUN sed -i 's/conn.SetReadDeadline(time.Now().Add(5 \* time.Minute))/conn.SetReadDeadline(time.Now().Add(5 \* time.Hour))/' /go/src/github.com/dedis/onet/websocket.go && \
    sed -i 's/const expirationTime = 1 \* time.Minute/const expirationTime = 1 \* time.Hour/' /go/src/github.com/dedis/onet/overlay.go && \
    cd /go/src/github.com/lca1/unlynx && \
    git checkout $UNLYNX_VERSION

# get sources and switch branch
WORKDIR /go/src/$MEDCO_UNLYNX_REPO
RUN git clone --depth 1 --branch $MEDCO_UNLYNX_VERSION https://$MEDCO_UNLYNX_REPO.git .

# get remaining dependencies, compile and install unlynxMedCo binary
# CGO_ENABLED=0 in order to be able to run from alpine
WORKDIR app
RUN go get -v -d ./... && \
    CGO_ENABLED=0 go build -v ./... && \
    CGO_ENABLED=0 go install -v ./...

# -------------------------------------------
FROM golang:1.11-alpine as prod

# run time environment variables
ENV NODE_IDX="0" \
    UNLYNX_DEBUG_LEVEL="1" \
    MEDCO_CONF_DIR="/medco-configuration"
ENV UNLYNX_BIN_EXPORT_PATH="$MEDCO_CONF_DIR/unlynxMedCo"

COPY --from=dev /go/bin/app /go/bin/unlynxMedCo
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh && \
    apk add --no-cache bash outils-md5

EXPOSE 2000 2001
VOLUME "$MEDCO_CONF_DIR"
WORKDIR "$MEDCO_CONF_DIR"
ENTRYPOINT docker-entrypoint.sh
