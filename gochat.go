package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"unsafe"
)

const (
	TYPE = "tcp"
	HOST = "localhost:8081"
)

func main() {
	server := flag.Bool("server", false, "start server")
	client := flag.Bool("client", false, "start client")
	message := Message{author: "Silly wizard", body: "Jesus, at least I tried."}
	encoded, _ := json.Marshal(message)
	var messageDecoded Message
	json.Unmarshal(encoded, &messageDecoded)
	fmt.Println(message.author+":", message.body)
	fmt.Println(unsafe.Sizeof(encoded))
	flag.Parse()

	if *server {
		createServer()
	}

	if *client {
		startClient()
	}
}
