[![Build Status](https://travis-ci.org/lca1/medco-loader.svg?branch=master)](https://travis-ci.org/LCA1/UnLynx) [![Go Report Card](https://goreportcard.com/badge/github.com/lca1/unlynx)](https://goreportcard.com/report/github.com/lca1/unlynx)

# MedCo Loader 
*medco-loader* is a small piece of software used to load data into MedCo, an operational system that makes sensitive medical data available for research in a simple, private and secure way. The current version offers two different loading models: (v0) loading of genomic data; and (v1) loading of encrypted i2b2 data. 
medco-loader is developed by lca1 (Laboratory for Communications and Applications in EPFL).  

## Documentation

* The medco-loader makes use of the [Advanced Crypto (kyber) library](https://github.com/dedis/kyber) and the [UnLynx library](https://github.com/lca1/unlynx) to encrypt the data. 
* For more information regarding the underlying crypto engine please refer to the stable versions of Kyber `github.com/dedis/kyber` and UnLynx `github.com/lca1/unlynx`
* To check the code organisation, have a look at [Layout](https://github.com/lca1/medco-loader/wiki/Layout)
* For more information on how to run the loader, go to [Running medco-loader](https://github.com/lca1/medco-loader/wiki/Running-medco-loader)

## Getting Started

To use the code of this repository you need to:

- Install [Golang](https://golang.org/doc/install)
- [Recommended] Install [IntelliJ IDEA](https://www.jetbrains.com/idea/) and the GO plugin
- Set [`$GOPATH`](https://golang.org/doc/code.html#GOPATH) to point to your workspace directory
- Add `$GOPATH/bin` to `$PATH`
- Git clone this repository to $GOPATH/src `git clone https://github.com/lca1/medco-loader.git` or...
- go get repository: `go get github.com/lca1/medco-loader`

## Version

The version in the `master`-branch is stable and has no incompatible changes.

## License

medco-loader is licensed under a End User Software License Agreement ('EULA') for non-commercial use. If you need more information, please contact us.

## Contact
You can contact any of the developers for more information or any other member of [lca1](http://lca.epfl.ch/people/lca1/):

* [Jean Louis Raisaro](https://github.com/JLRgithub) (Data Protection Specialist at CHUV) - jean.raisaro@epfl.ch
* [Juan Troncoso-Pastoriza](https://github.com/jrtroncoso) (Post-Doc) - juan.troncoso-pastoriza@epfl.ch
* [Mickaël Misbach](https://github.com/mickmis) (Software Engineer) - mickael.misbach@epfl.ch
* [Joao Andre Sa](https://github.com/JoaoAndreSa) (Software Engineer) - joao.gomesdesaesousa@epfl.ch
