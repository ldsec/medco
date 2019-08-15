FROM golang:1.12 as build

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

# compile and install medco-loader binary
# CGO_ENABLED=0 in order to be able to run from alpine
RUN CGO_ENABLED=0 go build -v ./... && \
    CGO_ENABLED=0 go install -v ./... && \
    mv /go/bin/app /go/bin/medco-loader

# -------------------------------------------
FROM golang:1.12-alpine as release

COPY --from=build /go/bin/medco-loader /go/bin/
RUN apk add --no-cache postgresql-client

# run-time environment
ENV LOG_LEVEL=5 \
    UNLYNX_GROUP_FILE_PATH=/medco-configuration/group.toml \
    UNLYNX_GROUP_FILE_IDX=0 \
    DB_HOST=localhost \
    DB_PORT=5432 \
    DB_NAME=i2b2medco \
    DB_USER=i2b2 \
    DB_PASSWORD=i2b2

ENTRYPOINT ["medco-loader"]
