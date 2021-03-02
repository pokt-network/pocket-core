# Development Guide
This guide contains useful knowledge for development.


## Setup a Development environment
In order to contribute first we need to setup our development environment

### Requirements
We need the following dependencies:
- [git](https://git-scm.com/)
- [go 1.14](https://golang.org/)
- [protobuffer compiler: protoc version 3.13.0](https://github.com/protocolbuffers/protobuf)
- A text editor of your choosing.

### Initial steps
First we need to clone our repository
```
$ git clone https://github.com/pokt-network/pocket-core.git
$ git checkout <Your brnach>
```

vendor dependencies
```
$ go mod vendor 
```

### Build
```
$ cd pocket-core/
$ go build -mod vendor -tags goleveldb -o /tmp/custom-pocket-build ./app/cmd/pocket_core/main.go
```

### Test
you may run all test witht he following command:
```
$ go test -p 1 ./...
```
This may take some time to complete, you can also run shorter tests with:
```
$ go test -p 1 -short ./...
```

or run module tests with:
```
$ go test -p ./x...
```

Mix and match on your needs.

## Contributing to Pocket Core
All contributions must come with an associated github issue.

Once you've created an issue on github feel free to clone, branch and develop your issue.

All PR's must be associated to an authored issue.

### Reaching out for help
You can reach out for by:
- [Writing issues on github](https://github.com/pokt-network/pocket-core/issues/new/choose)
- [Pocket ecosystem discord community](https://discord.com/invite/KRrqfd3tAK)

## Mocking Interfaces
Sometimes in order to test certain behaviours it's neccesary to use interfaces.

Our prime mocking candidate inteface is `Ctx` which denotes specific context expected (and unexpected).
Any update to the `Ctx` Interface would require an update of our mock struct inside the.

### Requirements
In order to create a mock structure we have the following dependency 
- [mockery](https://github.com/vektra/mockery)

This creates a struct that satisfies the Interface as well as use `github.com/stretchr/testify/mock` to provide values

### Usage
While usage instructions can be found on the official modkcery [documentation](https://github.com/vektra/mockery/blob/master/README.md) 

a quick usage for pocket core specific uses.
```
cd types/
mockery --name=Ctx
```

This creates a file inside a new directory 
```
 mocks/
   Ctx.go
```

Which contains an updated Ctx struct with ready to use.

Now we just need to move this Ctx struct onto `x/pocketcore/keeper/common_go` for usage

## Creating Proto Go types
As of RC-0.6.0 the adoption of protobuffers for encoding it's neccesary to update poto types which can be located within `proto/`

### Requirements
- protoc 3.13.0
- protoc-gen-gogo
- protoc-gen-grpc-gateway
- protoc-gen-swagger
- protoc-gen-go

### Usage
After installing neccesary third party tools.
Pocket provides an easy script for updating `.proto.pb` files

```
sh protoc/protocgen.sh
```
This will update all `proto.pb` files. 

