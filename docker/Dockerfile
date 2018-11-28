FROM golang:1.11

# environment variables
ENV MEDCO_LOADER_REPO="github.com/lca1/medco-loader" \
    MEDCO_LOADER_REPO_BRANCH="master" \
    BUILD_DIR="/opt/build-dir"

ENV MEDCO_LOADER_BIN_EXPORT_PATH="$BUILD_DIR/medco-loader"

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

RUN sed -i 's/conn.SetReadDeadline(time.Now().Add(5 \* time.Minute))/conn.SetReadDeadline(time.Now().Add(5 \* time.Hour))/' /go/src/github.com/dedis/onet/websocket.go && \
    sed -i 's/const expirationTime = 1 \* time.Minute/const expirationTime = 1 \* time.Hour/' /go/src/github.com/dedis/onet/overlay.go

# get sources and switch branch
WORKDIR /go/src/$MEDCO_LOADER_REPO
RUN git clone https://$MEDCO_LOADER_REPO.git . && \
    git checkout $MEDCO_LOADER_REPO_BRANCH

RUN go get -v -d ./...

# get remaining dependencies, compile and install medco binary, and rename it
WORKDIR app
RUN go install -v ./... && \
    mkdir -p $BUILD_DIR && \
    cp -a $(which app) $MEDCO_LOADER_BIN_EXPORT_PATH && \
    apt-get update && apt-get install -y postgresql-client-common postgresql-client && apt-get clean

ENTRYPOINT ["/opt/build-dir/medco-loader"]
