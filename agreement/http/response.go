package http

import (
	"fmt"

	"fly/common"
)

type (
	Response struct {
		Status  int
		Headers map[string]string
		Body    []byte
	}
)

func productionHttpResponse(req *Request) *Response {
	return &Response{
		Headers: req.Header,
		Status:  200,
	}
}

func (resp *Response) change2bytes() []byte {

	var stream []byte

	if resp.Status == 200 {
		stream = append(stream, []byte("HTTP/1.1 200 OK\r\n")...)
	} else {
		stream = append(stream, []byte("HTTP/1.1 404 Not Found\r\n")...)
	}

	resp.Headers["Content-Length"] = fmt.Sprintf("%v", len(resp.Body))

	for headerKey, headerValue := range resp.Headers {
		stream = append(stream, common.StringToBytes(headerKey)...)
		stream = append(stream, ':')
		stream = append(stream, common.StringToBytes(headerValue)...)
		stream = append(stream, []byte("\r\n")...)
	}

	stream = append(stream, []byte("\r\n")...)
	stream = append(stream, resp.Body...)
	return stream
}

func (resp *Response) SetStatus(code int) {
	resp.Status = code
}

func (resp *Response) Write(data []byte) {
	resp.Body = append(resp.Body, data...)
}
