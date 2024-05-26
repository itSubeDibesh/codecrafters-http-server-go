package main

import (
	"errors"
	"strings"
)

type Methods string

const (
	GET    Methods = "GET"
	POST   Methods = "POST"
	PUT    Methods = "PUT"
	DELETE Methods = "DELETE"
	PATCH  Methods = "PATCH"
)

type HTTPRequest struct {
	Method    string
	Protocol  string
	URI       string
	Host      string
	UserAgent string
	Accept    string
}

func ParseRequestRead(b []byte, n int) (*HTTPRequest, error) {
	if n == 0 || len(b) == 0 {
		return nil, errors.New("empty request")
	}
	bufferString := strings.Split(string(b[:n]), "\r\n")
	request := bufferString[0]
	r := &HTTPRequest{}
	r.Method = strings.Split(request, " ")[0]
	r.Protocol = strings.Split(request, " ")[2]
	r.URI = strings.Split(request, " ")[1]
	for _, v := range bufferString[1:] {
		if strings.Contains(v, "Host") {
			r.Host = strings.Split(v, " ")[1]
		}
		if strings.Contains(v, "User-Agent") {
			r.UserAgent = strings.Split(v, " ")[1]
		}
		if strings.Contains(v, "Accept") {
			r.Accept = strings.Split(v, " ")[1]
		}
	}
	return r, nil
}
