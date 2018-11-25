package main

import (
	"encoding/json"
	"fmt"
	"github.com/conformal/gotk3/gtk"
	// "github.com/nsallis/rpc-tester/callback"
	"net/rpc"
	"strings"
)

type Response struct {
	Payload string
}

type Request struct {
	Payload string
}

func main() {
	gtk.Init(nil)

	builder, err := gtk.BuilderNew()
	if err != nil {
		fmt.Println(err)
	}

	builder.AddFromFile("justlabel.ui")
	obj, err := builder.GetObject("window1")
	if err != nil {
		fmt.Println(err)
	}
	var window *gtk.Window
	if _, ok := obj.(*gtk.Window); !ok {
		panic("Window not found in .ui file")
	}
	window = obj.(*gtk.Window)
	window.Connect("destroy", gtk.MainQuit)

	sendButton, err := builder.GetObject("SendButton")
	sendButton.(*gtk.Button).Connect("clicked", func() {
		params := GetRequestParams(builder)
		SendRequest(params, builder)
		// TODO send request after getting params from GetRequestParams
	})

	window.ShowAll()
	gtk.Main()

	// callback.Init(builder)

}

type RequestParams struct {
	Address      string
	Port         string
	Command      string
	Payload      string
	FormatAsJson bool
}

// TODO error handling
func GetRequestParams(builder *gtk.Builder) RequestParams {
	// request address
	requestAddressObj, _ := builder.GetObject("RequestAddress")
	requestAddress, _ := requestAddressObj.(*gtk.Entry).GetText()

	// request port
	requestPortObj, _ := builder.GetObject("RequestPort")
	requestPort, _ := requestPortObj.(*gtk.Entry).GetText()

	// request command
	requestCommandObj, _ := builder.GetObject("RequestCommand")
	requestCommand, _ := requestCommandObj.(*gtk.Entry).GetText()

	// request text (payload)
	requestTextObj, _ := builder.GetObject("RequestText")
	requestTextBuffer, _ := requestTextObj.(*gtk.TextView).GetBuffer()
	start, end := requestTextBuffer.GetBounds()
	requestText, _ := requestTextBuffer.GetText(start, end, false)

	// format as json
	formatAsJsonObj, _ := builder.GetObject("FormatAsJson")
	formatAsJson := formatAsJsonObj.(*gtk.CheckButton).GetActive()

	params := RequestParams{
		Address:      requestAddress,
		Port:         requestPort,
		Command:      requestCommand,
		Payload:      requestText,
		FormatAsJson: formatAsJson,
	}
	return params
}

func SendRequest(params RequestParams, builder *gtk.Builder) {
	payload := params.Payload
	if params.FormatAsJson {
		payload = strings.Replace(payload, "\"", "\\\"", -1)
		payload = strings.Replace(payload, "\n", "\\\n", -1)
	}
	var (
		addr     = params.Address + ":" + params.Port
		request  = &Request{Payload: params.Payload}
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
	err = client.Call(params.Command, request, response)
	if err != nil {
		fmt.Printf("Could not call server: %v", err)
	}
	jsonResponse, err := json.Marshal(response.Payload)
	if err != nil {
		fmt.Println("Could not convert the response to json")
		fmt.Printf("Response: %v\n", response)
	}
	jsonText := strings.Replace(string(jsonResponse), "\\\"", "\"", -1)
	jsonText = strings.Replace(jsonText, ",\"", ",\n\"", -1)
	jsonText = strings.Replace(jsonText, "{\"", "{\n\"", -1)
	if err != nil {
		fmt.Sprintln("error pretty printing: %e", err)
	}
	fmt.Println("jsonText")
	fmt.Println(jsonText)

	jsonText = strings.Replace(string(jsonText), ":{", ":\n\t{", -1)

	responseTextObj, _ := builder.GetObject("ResponseText")
	responseTextBox := responseTextObj.(*gtk.TextView)
	responseTextBuffer, _ := responseTextBox.GetBuffer()
	responseTextBuffer.SetText(string(jsonText))
	fmt.Println("Response: " + string(jsonText))

}
