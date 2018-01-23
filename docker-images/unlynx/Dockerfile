FROM golang:1.8

# environment variables
ENV UNLYNX_REPO="github.com/lca1/unlynx" \
    UNLYNX_REPO_BRANCH="medco-deployment" \
    CONF_DIR="/opt/medco-configuration"

# get maximum of dependencies and cache them
RUN go get -v -d gopkg.in/dedis/onet.v1/... && \
    go get -v -d gopkg.in/dedis/crypto.v0/... && \
    go get -v -d gopkg.in/urfave/cli.v1/... && \
    go get -v -d github.com/Knetic/govaluate/... && \
    go get -v -d github.com/btcsuite/goleveldb/... && \
    go get -v -d github.com/r0fls/gostats/... && \
    go get -v -d github.com/fanliao/go-concurrentMap/... && \
    go get -v -d github.com/lib/pq/...

# get sources and switch branch
WORKDIR /go/src/$UNLYNX_REPO
RUN git clone https://$UNLYNX_REPO.git . && \
    git checkout $UNLYNX_REPO_BRANCH

# get remaining dependencies, compile and install unlynxMedCo binary, and make it available in the shared folder
WORKDIR app/unlynxMedCo
RUN go get -v -d ./... && \
    go install -v ./...

# add script that copies the binary in the configuration folder
COPY copy-unlynx-binary.sh /

# conf and run
ARG NODE_IDX_ARG="0"
ARG UNLYNX_DEBUG_LEVEL_ARG="1"
ENV NODE_IDX="$NODE_IDX_ARG" \
    UNLYNX_BIN_EXPORT_PATH="$CONF_DIR/unlynxMedCo" \
    UNLYNX_DEBUG_LEVEL="$UNLYNX_DEBUG_LEVEL_ARG"
ENV UNLYNX_KEY_FILE_PATH="$CONF_DIR/srv$NODE_IDX-private.toml" \
    UNLYNX_DDT_SECRETS_FILE_PATH="$CONF_DIR/srv$NODE_IDX-ddtsecrets.toml"

EXPOSE 2000 2001
VOLUME "$CONF_DIR"
WORKDIR "$CONF_DIR"
ENTRYPOINT ["sh", "-c", "sh /copy-unlynx-binary.sh && $UNLYNX_BIN_EXPORT_PATH -d $UNLYNX_DEBUG_LEVEL server -c $UNLYNX_KEY_FILE_PATH"]
