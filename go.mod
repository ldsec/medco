module github.com/CHUV-DS/medco

replace github.com/ldsec/medco => ./

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/fanliao/go-concurrentMap v0.0.0-20141114143905-7d2d7a5ea67b
	github.com/go-openapi/errors v0.20.1
	github.com/go-openapi/loads v0.21.0
	github.com/go-openapi/runtime v0.21.0
	github.com/go-openapi/spec v0.20.4
	github.com/go-openapi/strfmt v0.21.1
	github.com/go-openapi/swag v0.19.15
	github.com/go-openapi/validate v0.20.3
	github.com/jessevdk/go-flags v1.5.0
	github.com/ldsec/medco v0.0.0-00010101000000-000000000000
	github.com/ldsec/unlynx v1.4.3
	github.com/lestrrat-go/jwx v1.2.13
	github.com/lib/pq v1.10.4
	github.com/pkg/errors v0.9.1
	github.com/r0fls/gostats v0.0.0-20180711082619-e793b1fda35c
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.5
	go.dedis.ch/kyber/v3 v3.0.13
	go.dedis.ch/onet/v3 v3.2.10
	golang.org/x/net v0.0.0-20211209124913-491a49abca63
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

go 1.15
