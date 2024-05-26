package main

type HTTPResponse string

const (
	OK                  HTTPResponse = "HTTP/1.1 200 OK\r\n\r\n"
	NotFound            HTTPResponse = "HTTP/1.1 404 Not Found\r\n\r\n"
	Created             HTTPResponse = "HTTP/1.1 201 Created\r\n\r\n"
	InternalServerError HTTPResponse = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
)
