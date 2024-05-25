package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Initializing server...")

	const PORT = 4221

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

		handelClientConnection(connection)
	}
}

func handelClientConnection(connection net.Conn) {
	defer connection.Close()
	buffer := make([]byte, 1024)
	_, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}

	_, err = connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		return
	}

}
