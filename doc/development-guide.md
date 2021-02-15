# Development Guide
This guide contains useful knowledge for development.

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

