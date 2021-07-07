# cmp-ecdsa

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A Go implementation of the "CMP" protocol by [Canetti et al.](https://eprint.iacr.org/2021/060) for threshold ECDSA signing, with some additions to improve its practical reliability, including the "echo broadcast" from [Goldwasser and Lindell](https://doi.org/10.1007/s00145-005-0319-z).  We documented these in [threshold_protocol.pdf](threshold_protocol.pdf).

Our implementation supports ECDSA with secp256k1.

## Usage

### Sessions

Protocol data is stored in a `session.Session` struct.

### Keygen

```go
baseSession, err := session.NewKeygenSession([]party.ID{"a","b","c"}, 2, "a") (*Keygen, error) 
if err != nil {
	//handle err
}
p, err := protocols.NewProtocol(baseSession, refresh.Create)
if err != nil {
    //handle err
}
inChan, outChan, errChan := p.Start()
// TODO handle message passing 
refreshedSession, err := p.GetSession()
if err != nil {
    //handle err
}
```
### Refresh

```go
var oldSession session.Session
if err != nil {
	//handle err
}
p, err := protocols.NewProtocol(baseSession, refresh.Create)
if err != nil {
    //handle err
}
inChan, outChan, errChan := p.Start()
// TODO handle message passing 
refreshedSession, err := p.GetSession()
if err != nil {
    //handle err
}

```
### Sign

For now, we only implement the 4 round interactive signing protocol. 
```go

var (
	refreshedSession session.Session
	message []byte
	signers = []party.ID
)
signSession, err :=  NewSignSession(refreshedSession, signerIDs, message)
if err != nil {
	//handle err
}
p, err := protocols.NewProtocol(signSession, sign.Create)
if err != nil {
    //handle err
}
inChan, outChan, errChan := p.Start()
// TODO handle message passing 
sig, err := p.GetSignature()
if err != nil {
    //handle err
}
```
### Network

The messages returned by the protocol can be transmitted through an authenticated point-to-point network.

## Build

`cmp-ecdsa` requires a custom version of `gogo` which enables the use of `*big.Int` in protobufs. 
This version can be compiled by applying the path from [trasc/casttypewith](https://github.com/trasc/protobuf)
It can be installed using the following shell commands:

```shell
git clone https://github.com/gogo/protobuf $GOPATH/src/github.com/gogo/protobuf
cd $GOPATH/src/github.com/gogo/protobuf
git remote add trasc https://github.com/trasc/protobuf.git
git fetch trasc
git merge trasc/casttypewith
cd protoc-gen-gogofaster
go build
cp protoc-gen-gogofaster $GOPATH/bin
```

Once installed, running `make` in the root will regenerate all `.proto` files.


## Intellectual property

This code is copyright (c) Adrian Hamelink and Taurus SA, 2021, and under Apache 2.0 license.




