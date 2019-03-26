FROM medco/unlynx:v1.2b-alpha as build

# environment variables
ENV MEDCO_LOADER_REPO="github.com/lca1/medco-loader"

ENV GO111MODULE=on

# adjust timeout values
RUN sed -i 's/conn.SetReadDeadline(time.Now().Add(5 \* time.Minute))/conn.SetReadDeadline(time.Now().Add(5 \* time.Hour))/' /go/src/github.com/dedis/onet/websocket.go && \
    sed -i 's/const expirationTime = 1 \* time.Minute/const expirationTime = 1 \* time.Hour/' /go/src/github.com/dedis/onet/overlay.go

COPY ./ /go/src/$MEDCO_LOADER_REPO

# get remaining dependencies, compile and install medco-loader binary
# CGO_ENABLED=0 in order to be able to run from alpine
WORKDIR /go/src/$MEDCO_LOADER_REPO
RUN go get -v -d ./... && \
    CGO_ENABLED=0 go build -v ./... && \
    CGO_ENABLED=0 go install -v ./... && \
    mv /go/bin/app /go/bin/medco-loader && \
    apt-get update && apt-get install -y postgresql-client-common postgresql-client && apt-get clean


# -------------------------------------------
FROM golang:1.11-alpine as release

COPY --from=build /go/bin/medco-loader /go/bin/
RUN apk add --no-cache postgresql-client
ENTRYPOINT ["medco-loader"]
