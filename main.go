package main

import (
	"flag"
	"fmt"
	"net/rpc"
  "encoding/json"
)

type Response struct {
	Payload string
}

type Request struct {
	Payload     string
}

func main() {
	portUsage := "Which port the rpc server is running on"
	commandUsage := "Which command to invoke"
	payloadUsage := "Payload to send"
	var port string
	var command string
	var payload string
	flag.StringVar(&port, "port", "3333", portUsage)
	flag.StringVar(&port, "p", "3333", portUsage)

	flag.StringVar(&command, "command", "Endpoints.Handshake", commandUsage)
	flag.StringVar(&command, "c", "Endpoints.Handshake", commandUsage)

	flag.StringVar(&payload, "payload", "", payloadUsage)
	flag.StringVar(&payload, "P", "", payloadUsage)

	flag.Parse()

	var (
		addr     = "localhost:" + port
		request  = &Request{Payload: payload}
		response = new(Response)
	)
	// Establish the connection to the adddress of the
	// RPC server
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Error dialing server: %v\n", err)
		fmt.Println("Is the server running?")
		return
	}
	defer client.Close()

	// Perform a procedure call (core.HandlerName == Handler.Execute)
	// with the Request as specified and a pointer to a response
	// to have our response back.
	err = client.Call(command, request, response)
	if err != nil {
		fmt.Printf("Could not call server: %v", err)
	}
  jsonResponse, err := json.Marshal(response)
  if err != nil {
    fmt.Println("Could not convert the response to json")
    fmt.Printf("Response: %v\n", response)
  }
	fmt.Println("Response: " + string(jsonResponse))
}
