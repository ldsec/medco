[![Build Status](https://travis-ci.org/lca1/medco-unlynx.svg?branch=master)](https://travis-ci.org/lca1/medco-unlynx) 
[![Go Report Card](https://goreportcard.com/badge/github.com/lca1/medco-unlynx)](https://goreportcard.com/report/github.com/lca1/medco-unlynx) 
[![Coverage Status](https://coveralls.io/repos/github/lca1/medco-unlynx/badge.svg?branch=master)](https://coveralls.io/github/lca1/medco-unlynx?branch=master)

## Documentation
MedCo documentation is centralized on the following website: 
[MedCo Unlynx](https://medco.epfl.ch/documentation/developer/components/medco-unlynx.html).

## Version

We have a development and a stable version. The `master`-branch in `github.com/lca1/medco-unlynx` is the development version that works but can have incompatible changes.

Use one of the latest tags `v0.1.1d` that are stable and have no incompatible changes.

**Very Important!!** 

Due to the current changes being made to [onet](https://github.com/dedis/onet) and [kyber](https://github.com/dedis/kyber) (release of v3) you must revert back to previous commits for these two libraries if you want medco-unlynx to work. This will change in the near future. 

```bash
cd $GOPATH/src/dedis/onet/
git checkout 5796104343ef247e2eed58e573f68c566db2136f

cd $GOPATH/src/dedis/kyber/
git checkout f55fec5463cda138dfc7ff15e4091d12c4ddcbfe
```

## License
*medco-unlynx* is licensed under a End User Software License Agreement ('EULA') for non-commercial use. 
If you need more information, please contact us.
