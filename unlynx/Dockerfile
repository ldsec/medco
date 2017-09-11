FROM golang:1.8

# environment variables
ENV UNLYNX_REPO="github.com/lca1/unlynx" \
    UNLYNX_REPO_BRANCH="newI2B2"

# get maximum of dependencies and cache them
RUN go get -v -d gopkg.in/dedis/onet.v1/... && \
    go get -v -d gopkg.in/dedis/crypto.v0/... && \
    go get -v -d gopkg.in/urfave/cli.v1/... && \
    go get -v -d github.com/Knetic/govaluate/... && \
    go get -v -d github.com/btcsuite/goleveldb/... && \
    go get -v -d github.com/r0fls/gostats/... && \
    go get -v -d github.com/fanliao/go-concurrentMap/...

# get sources and switch branch
WORKDIR /go/src/$UNLYNX_REPO
RUN git clone https://$UNLYNX_REPO.git . && \
    git checkout $UNLYNX_REPO_BRANCH

# get remaining dependencies, compile and install unlynxI2b2 binary, and make it available in the shared folder
WORKDIR app/i2b2
RUN go get -v -d ./... && \
    go install -v ./...

# conf and run
ARG UNLYNX_KEY_FILE_ARG="private.toml"
ENV UNLYNX_KEY_FILE="$UNLYNX_KEY_FILE_ARG" \
    UNLYNX_BIN_EXPORT_PATH="/opt/configuration/unlynxI2b2"
EXPOSE 2000 2001
VOLUME /opt/configuration
WORKDIR /opt/configuration
ENTRYPOINT ["sh", "-c", "rm -f $UNLYNX_BIN_EXPORT_PATH && cp -a $(which i2b2) $UNLYNX_BIN_EXPORT_PATH && chmod 777 $UNLYNX_BIN_EXPORT_PATH && $UNLYNX_BIN_EXPORT_PATH -d 1 server -c /opt/configuration/keys/$UNLYNX_KEY_FILE"]
# TODO switch for prod:
# TODO: debug option to modify
# TODO: logs to put in unlynx folder (now: stdin)
# TODO: key generation handling (docker exec? or get the unlynx volume?)
