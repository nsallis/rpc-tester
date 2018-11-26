# rpc-tester
Simple rpc server tester with options for address, port, command, and payload.
Request is sent as {Payload: PAYLOAD}
If the Send As Json checkbox is checked, the request payload is escaped before sending (you can write your payload as `{"foo":"bar"}`, and unmarshaling on the rpc side will result in the correct result.), and the response is formatted as json (though not perfectly. Suggestions on how to better format json strings are welcome!).
