package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/google/uuid"
)

const (
	CONNECTION_TYPE = "tcp"
	HOST            = "localhost"
	PORT            = "8081"
)

type Server struct {
	listener net.Listener
	messages *chan *Message
	clients  map[uuid.UUID]*Client
}

type Client struct {
	conn net.Conn
	name string
	id   uuid.UUID
}

type Message struct {
	Body string
	From string
}

func main() {
	serverFlag := flag.Bool("server", false, "Starts a server.")
	clientFlag := flag.Bool("client", false, "Starts a client")
	flag.Parse()
	if *serverFlag {
		startServer()
	} else if *clientFlag {
		startClient()
	} else {
		flag.Usage()
	}
}

// startServer sets up a listener to listen for new clients
func startServer() {
	fmt.Println("Starting Server")
	listener, err := net.Listen(CONNECTION_TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error starting server.")
		os.Exit(1)
	}
	defer listener.Close()
	// Create server object
	messages := make(chan *Message)
	server := Server{listener: listener, messages: &messages, clients: make(map[uuid.UUID]*Client)}

	// Begin accepting new clients
	go acceptClients(&server)
	awaitMessagesServer(&server)
}

// acceptClients accepts new clients connecting to the server then handles them
func acceptClients(server *Server) {
	for {
		// Accept new client
		conn, err := server.listener.Accept()
		if err != nil {
			fmt.Println("Error connecting to client:", err.Error())
			continue
		}
		// Create client object and add to server slice
		client := Client{conn: conn, name: fmt.Sprint("Client", len(server.clients)), id: uuid.New()}
		server.clients[client.id] = &client
		go listenToClient(server, &client)
		fmt.Println(client.name, "connected.")
	}
}

// listenToClient waits for the client to send a message, then handles that message
func listenToClient(server *Server, client *Client) {
	buffer := make([]byte, 1024)
	for {
		n, err := client.conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println(client.name, "disconnected.")
				server.clients[client.id].conn.Close()
				delete(server.clients, client.id)
				break
			}
			fmt.Println("Error reading message from", client.name, ":", err.Error())
			continue
		}
		message := Message{From: client.name, Body: string(buffer[:n])}
		// err = json.Unmarshal(buffer[:n], &message)
		// if err != nil {
		// 	fmt.Println("Error decoding message from", client.name, ":", err.Error())
		// 	continue
		// }
		*server.messages <- &message
	}
}

// awaitMessagesServer awaits messages sent to the server and processes them
func awaitMessagesServer(server *Server) {
	for {
		message := <-*server.messages
		fmt.Println(message.From+":", message.Body)
		encodedMessage, err := json.Marshal(*message)
		if err != nil {
			fmt.Println("Error marshling message:", err.Error())
			continue
		}
		for _, client := range server.clients {
			if client.name != message.From {
				client.conn.Write(encodedMessage)
			}
		}
	}
}

// startClient starts the client and begins the handling processes
func startClient() {
	conn, err := net.Dial(CONNECTION_TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error connecting to server:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	go awaitMessagesClient(&conn)
	awaitInputClient(&conn)
}

// awaitMessagesClient reads and displays messages recieved from the server
func awaitMessagesClient(conn *net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := (*conn).Read(buffer)
		if err != nil {
			fmt.Println("Error reading message from server:", err.Error())
			continue
		}
		var message Message
		err = json.Unmarshal(buffer[:n], &message)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err.Error())
			continue
		}
		// fmt.Println(message.From+":", message.Body)
		fmt.Printf("\r" + message.From + ":" + message.Body + "\n")
	}
}

// awaitInputClient awaits input from the user, then sends that text to the server
func awaitInputClient(conn *net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		(*conn).Write([]byte(scanner.Text()))
	}
}
