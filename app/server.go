package main

import (
	"fmt"
	"net"
	"os"
)

const (
	MAX_BUFFER_SIZE = 1024
	PORT            = 4221
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Initializing server...")

	// Used for inbound connections -> We only use it as we are server not client
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", PORT))
	if err != nil {
		fmt.Printf("Failed to bind to port %v\n", PORT)
		os.Exit(1)
	}

	defer listener.Close()
	fmt.Printf("Listening on port %v\n", PORT)

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		go handleClientConnection(connection)
	}
}

func handleClientConnection(connection net.Conn) {
	defer connection.Close()

	var response HTTPResponse

	buffer := make([]byte, MAX_BUFFER_SIZE)
	requestTimeout, err := connection.Read(buffer)

	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}

	request, err := ParseRequestRead(buffer, requestTimeout)
	if err != nil {
		fmt.Println("Error parsing request: ", err.Error())
		return
	}

	if request.URI == "/" {
		response = OK
	} else {
		response = NotFound
	}
	if _, err = connection.Write([]byte(response)); err != nil {
		fmt.Println("Error writing response: ", err.Error())
		return
	}

}
