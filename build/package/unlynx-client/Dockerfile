FROM golang:1.14 as build
ARG MEDCO_VERSION=dev

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

# compile and install all binaries
RUN CGO_ENABLED=0 go install -v -ldflags "-X github.com/ldsec/medco.Version=${MEDCO_VERSION}" ./...

# -------------------------------------------
FROM golang:1.14-alpine as release

COPY build/package/unlynx-client/docker-entrypoint.sh /usr/local/bin/
RUN apk update && apk add bash && rm -rf /var/cache/apk/* && \
    chmod a+x /usr/local/bin/docker-entrypoint.sh
COPY --from=build /go/bin /go/bin

ENV UNLYNX_GROUP_FILE_PATH=/medco-configuration/group.toml \
    MEDCO_NODE_IDX=0 \
    UNLYNX_DEBUG_LEVEL="1"

VOLUME "/medco-configuration"
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["ls -1 /go/bin"]
