# rpc-tester
Simple rpc server tester with options for port, command, and payload.
Request is sent as {Payload: PAYLOAD}
the response is parsed as json, or printed in go format if that fails.

Example:

```
go build main.go test
./test -p 3333 -c "Endpoints.Handshake" -P "some payload"
```
