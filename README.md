## Go json-rpc protobuf plugin
JSON-RPC? gRPC? jRPC!

jRPC is a simple plugin that generates http endpoints from a proto schema. It is useful when your client does not have access to gRPC endpoints

Try it:
`git clone https://github.com/danielvladco/jrpc.git && cd ./jrpc/example && go run main.go` 

Test it:
`curl --data '{ "param1": "str", "param2": 1 }'  http://localhost:8080/Service/Endpoint1`

Generate command:
`protoc --go_out=plugins=grpc,paths=source_relative:. --jrpc_out=paths=source_relative:. ./example/pb/*.proto`

It's as easy as this! try it out and save yourself lots of time writing boilerplate code and debugging potential bugs.

See `/example`  for more documentation

## Features
- JSON Request decoding
- JSON response encoding
- Generates endpoints by the following formula: `/` + service name + `/` + method name. ex: `/Service/Endpoint1`
- gRPC status codes to http status codes
- (not supported) client streaming
- (not supported) server streaming

## License:

jRPC is released under the Apache 2.0 license. See the LICENSE file for details.
