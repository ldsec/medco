module github.com/ldsec/medco-unlynx

go 1.13

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/btcsuite/goleveldb v1.0.0
	github.com/fanliao/go-concurrentMap v0.0.0-20141114143905-7d2d7a5ea67b
	github.com/ldsec/unlynx v1.4.0
	github.com/smartystreets/goconvey v1.6.4 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.22.3
	go.dedis.ch/kyber/v3 v3.0.12
	go.dedis.ch/onet/v3 v3.1.1
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/net v0.0.0-20191126235420-ef20fe5d7933 // indirect
	golang.org/x/sys v0.0.0-20200316230553-a7d97aace0b0 // indirect
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace github.com/ldsec/unlynx => ../unlynx
