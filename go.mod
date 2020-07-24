module github.com/ldsec/medco-connector

replace github.com/ldsec/medco-unlynx => github.com/ldsec/medco-unlynx v0.3.2-0.20200414130428-b1239ae61d90

require (
	github.com/go-openapi/errors v0.19.6
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/runtime v0.19.20
	github.com/go-openapi/spec v0.19.8
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/go-openapi/validate v0.19.10
	github.com/go-swagger/go-swagger v0.25.0 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/ldsec/medco-loader v1.0.0
	github.com/ldsec/medco-unlynx v1.0.0
	github.com/ldsec/unlynx v1.4.1
	github.com/lestrrat-go/jwx v0.9.0
	github.com/lib/pq v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/r0fls/gostats v0.0.0-20180711082619-e793b1fda35c
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/urfave/cli v1.22.3
	github.com/vektah/gqlparser v1.1.2 // indirect
	go.dedis.ch/onet/v3 v3.2.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/tools/gopls v0.4.3 // indirect
)

go 1.13
