FROM golang:1.13 as build

COPY ./ /src
WORKDIR /src

# get dependencies
RUN go get -v -d ./...

# hack to adjust timeout values in onet
RUN cd /go/pkg/mod/go.dedis.ch/onet && \
    chmod u+w -R . && \
    find . -maxdepth 2 -name overlay.go | xargs sed -i \
        's/const expirationTime = 1 \* time.Minute/const expirationTime = 1 \* time.Hour/' && \
    find . -maxdepth 2 -name websocket.go | xargs sed -i \
        's/conn.SetReadDeadline(time.Now().Add(5 \* time.Minute))/conn.SetReadDeadline(time.Now().Add(5 \* time.Hour))/'

# compile and install medco-unlynx binary
# CGO_ENABLED=0 in order to be able to run from alpine
RUN CGO_ENABLED=0 go build -v ./... && \
    CGO_ENABLED=0 go install -v ./... && \
    mv /go/bin/app /go/bin/medco-unlynx

# -------------------------------------------
FROM golang:1.13-alpine as release

# run time environment variables
ENV NODE_IDX="0" \
    UNLYNX_DEBUG_LEVEL="1" \
    MEDCO_CONF_DIR="/medco-configuration"

COPY --from=build /go/bin/medco-unlynx /go/bin/
COPY deployment/docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh && \
    apk add --no-cache bash

VOLUME "$MEDCO_CONF_DIR"
ENTRYPOINT ["docker-entrypoint.sh"]
