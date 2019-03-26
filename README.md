[![Build Status](https://travis-ci.org/lca1/medco-loader.svg?branch=master)](https://travis-ci.org/lca1/medco-loader) 
[![Go Report Card](https://goreportcard.com/badge/github.com/lca1/medco-loader)](https://goreportcard.com/report/github.com/lca1/medco-loader) 
[![Coverage Status](https://coveralls.io/repos/github/lca1/medco-loader/badge.svg?branch=master)](https://coveralls.io/github/lca1/medco-loader?branch=master)

## Documentation
MedCo documentation is centralized on the following website: 
[MedCo Loader](https://medco.epfl.ch/documentation/developer/components/medco-loader.html).

## Version

We have a development and a stable version. The `master`-branch in `github.com/lca1/medco-loader` is the development version that works but can have incompatible changes.

Use one of the latest tags `v0.1.1a` that are stable and have no incompatible changes.

**Very Important!!** 

Due to the current changes being made to [onet](https://go.dedis.ch/onet/v3) and [kyber](https://go.dedis.ch/kyber/v3) (release of v3) you must revert back to previous commits for these two libraries if you want medco-loader to work. This will change in the near future. 

```bash
cd $GOPATH/src/dedis/onet/
git checkout 5796104343ef247e2eed58e573f68c566db2136f

cd $GOPATH/src/dedis/kyber/
git checkout f55fec5463cda138dfc7ff15e4091d12c4ddcbfe
```

## License
*medco-loader* is licensed under a End User Software License Agreement ('EULA') for non-commercial use.
If you need more information, please contact us.
