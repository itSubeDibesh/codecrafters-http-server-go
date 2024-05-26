package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
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

	response = handleRequest(request)

	if _, err = connection.Write([]byte(response)); err != nil {
		fmt.Println("Error writing response: ", err.Error())
		return
	}
}

func handleRequest(request *HTTPRequest) HTTPResponse {
	if request.URI == "/" {
		return OK
	} else if isValidEchoURI(request.URI) {
		return handleEchoRequest(request)
	} else if strings.HasPrefix(request.URI, "/user-agent") {
		return handleUserAgentRequest(request)
	} else if isValidFileURI(request.URI) {
		if request.Method == string(GET) {
			return handleFileRequest(request)
		}
		if request.Method == string(POST) {
			return handelFileUploadRequest(request)
		}
		return NotFound
	} else {
		return NotFound
	}
}

func isValidEchoURI(uri string) bool {
	parts := strings.Split(uri, "/")
	return len(parts) == 3 && parts[1] == "echo" && strings.HasPrefix(uri, "/echo")
}

func handleEchoRequest(request *HTTPRequest) HTTPResponse {
	pathParams := strings.Split(request.URI, "/")[2]
	value, ok := request.Headers["Accept-Encoding"]
	if ok {
		if strings.Contains(strings.ToLower(strings.TrimSpace(value)), "gzip") {
			buffer := new(bytes.Buffer)
			gzWrite := gzip.NewWriter(buffer)
			gzWrite.Write([]byte(pathParams))
			gzWrite.Close()
			return getSuccessResponse(buffer.String(), "gzip")
		}
		return getSuccessResponse(pathParams, "")
	}
	return getSuccessResponse(pathParams, "")
}

func handleUserAgentRequest(request *HTTPRequest) HTTPResponse {
	return getSuccessResponse(request.UserAgent, "")
}

func isValidFileURI(uri string) bool {
	parts := strings.Split(uri, "/")
	return len(parts) == 3 && parts[1] == "files" && strings.HasPrefix(uri, "/files")
}

func handleFileRequest(request *HTTPRequest) HTTPResponse {
	file := strings.Split(request.URI, "/")[2]
	dir := os.Args[2]
	if _, err := os.Stat(filepath.Join(dir, file)); os.IsNotExist(err) {
		return NotFound
	} else {
		data, err := os.ReadFile(filepath.Join(dir, file))
		if err != nil {
			return NotFound
		} else {
			return HTTPResponse(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%v\r\n", len(data), string(data)))
		}
	}
}

func getSuccessResponse(content string, encoding string) HTTPResponse {
	if len(encoding) != 0 {
		return HTTPResponse(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nContent-Encoding: %s\r\n\r\n%s\r\n", len(content), encoding, content))
	}
	return HTTPResponse(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s\r\n", len(content), content))
}

func handelFileUploadRequest(request *HTTPRequest) HTTPResponse {
	file := strings.Split(request.URI, "/")[2]
	dir := os.Args[2]
	data := []byte(bytes.Trim([]byte(request.Body), "\x00"))

	if err := os.WriteFile(filepath.Join(dir, file), data, 0644); err != nil {
		return InternalServerError
	}
	return Created
}
