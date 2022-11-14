package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Server struct {
	listener net.Listener
	clients  []Client
	messages chan Message
}

func createServer() {
	server := initServer()
	defer shutdownServer(server)
	startServer(server)
}

func shutdownServer(server *Server) {
	fmt.Println("Shutting down server.")
	server.listener.Close()
	for _, client := range server.clients {
		client.connection.Write([]byte("shutdown"))
	}
}

/**
* Initializes server
* Return: *Server - Pointer to server struct
 */
func initServer() *Server {
	fmt.Println("Creating server.")
	listener, err := net.Listen(TYPE, HOST)
	if err != nil {
		exitWithError("Failed to start listener", err)
	}
	fmt.Println("Sever created.")

	server := Server{listener, make([]Client, 0), make(chan Message)}
	return &server
}

/**
* Starts listening and accepting clients.  Calls handleClient on each new connection
 */
func startServer(server *Server) {
	fmt.Println("Listening on server.")
	go handleMessages(server)
	for {
		connection, err := server.listener.Accept()
		if err != nil {
			warn("Could not accept connection", err)
		}
		fmt.Println("Listening to new client")
		server.clients = append(server.clients, Client{connection: connection})
		go handleClient(connection, server)
	}
}

/**
* Begins listening to a client for new messages
 */
func handleClient(connection net.Conn, server *Server) {
	for {
		buffer := make([]byte, 1024)
		_, err := connection.Read(buffer)
		if err != nil {
			warn("Error reading from client", err)
			continue
		}
		var message Message
		err = json.Unmarshal(buffer, &message)
		if err != nil {
			warn("Error unmarshling message", err)
			continue
		}
		fmt.Println(message.author+":", message.body)
		server.messages <- message
	}
}

/**
* Waits for a new message in the message channel, then send it to each client
 */
func handleMessages(server *Server) {
	for {
		outgoing, err := json.Marshal(<-server.messages)
		if err != nil {
			warn("Error marshaling message", err)
			continue
		}
		fmt.Println("Handling message")
		for _, client := range server.clients {
			_, err := client.connection.Write(outgoing)
			if err != nil {
				warn("Error sending message to "+client.name, err)
			}
			fmt.Println("Sent to:", client.name)
		}
	}
}

func exitWithError(message string, err error) {
	fmt.Println("[X]", message, ":", err.Error())
	os.Exit(1)
}

func warn(message string, err error) {
	fmt.Println("[!]", message, ":", err.Error())
}
