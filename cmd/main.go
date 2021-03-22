package main

import (
	"flag"
	"log"

	"../pkg/UI"
	"../pkg/client"
)

func main() {
	address := flag.String("server", "0.0.0.0:3333", "Which server to connect to")

	flag.Parse()

	client := client.NewClient()
	err := client.Dial(*address)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	// start the client to listen for incoming message
	go client.Start()

	UI.StartUi(client)
}
