FROM golang:1.10.3

# environment variables
ENV UNLYNX_REPO="github.com/lca1/medco" \
    UNLYNX_REPO_BRANCH="master" \
    BUILD_DIR="/opt/build-dir"

ENV UNLYNX_BIN_EXPORT_PATH="$BUILD_DIR/unlynx"

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

# get sources and switch branch
WORKDIR /go/src/$UNLYNX_REPO
RUN git clone https://$UNLYNX_REPO.git . && \
    git checkout $UNLYNX_REPO_BRANCH

# get remaining dependencies, compile and install medco binary, and make it available in the shared folder
WORKDIR app
RUN go get -v -d ./... && \
    go install -v ./...

COPY copy-unlynx-binary-dev.sh /
ENTRYPOINT ["sh", "-c", "/copy-unlynx-binary-dev.sh"]