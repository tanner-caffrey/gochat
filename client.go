package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Client struct {
	connection net.Conn
	name       string
}

type Message struct {
	body   string
	author string
}

func startClient() {
	client := createClient()
	defer stopClient(client)
	go acceptInput(client)
	listenToServer(client)
}

func stopClient(client *Client) {
	client.connection.Close()
	fmt.Println("Shutting down client.")
}

func createClient() *Client {
	fmt.Printf("Enter name:")
	var name string
	fmt.Scanln(&name)
	fmt.Println("Creating client connection")
	connection, err := net.Dial(TYPE, HOST)
	if err != nil {
		exitWithError("Error connecting to server", err)
	}
	fmt.Println("Connected to", HOST)
	client := Client{connection: connection, name: name}
	return &client
}

func acceptInput(client *Client) {
	var outgoingMessage string
	for {
		// fmt.Printf(">")
		fmt.Scanln(&outgoingMessage)
		message := Message{body: outgoingMessage[:len(outgoingMessage)-1], author: client.name}
		bytes, err := json.Marshal(message)
		if err != nil {
			warn("Error encoding message", err)
		}
		client.connection.Write(bytes)
	}
}

func listenToServer(client *Client) {
	fmt.Println("Listening to server")
	buffer := make([]byte, 1024)
	for {
		n, err := client.connection.Read(buffer)
		if err != nil {
			warn("Error reading message from server", err)
			continue
		}
		var message Message
		err = json.Unmarshal(buffer[:n], &message)
		if err != nil {
			warn("Error unmarshling message", err)
		}
		if string(buffer[:n]) == "shutdown" {
			break
		}
		fmt.Println(message.author+":", message.body)
	}
}
