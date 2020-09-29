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

# compile and install medco-loader binary
# CGO_ENABLED=0 in order to be able to run from alpine
RUN CGO_ENABLED=0 go build -v ./... && \
    CGO_ENABLED=0 go install -v ./... && \
    mv /go/bin/app /go/bin/medco-loader

# -------------------------------------------
FROM golang:1.13-alpine as release

COPY --from=build /go/bin/medco-loader /go/bin/
RUN apk add --no-cache postgresql-client

# run-time environment
ENV LOG_LEVEL=5 \
    UNLYNX_GROUP_FILE_PATH=/medco-configuration/group.toml \
    UNLYNX_GROUP_FILE_IDX=0 \
    I2B2_DB_HOST=localhost \
    I2B2_DB_PORT=5432 \
    I2B2_DB_NAME=i2b2medco \
    I2B2_DB_USER=i2b2 \
    I2B2_DB_PASSWORD=i2b2 \
    GA_DB_HOST=localhost \
    GA_DB_PORT=5432 \
    GA_DB_NAME=gamedcosrv \
    GA_DB_USER=genomicannotations \
    GA_DB_PASSWORD=genomicannotations

ENTRYPOINT ["medco-loader"]
