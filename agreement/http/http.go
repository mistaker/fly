package http

import (
	"fly/server"
	"fly/socket"
)

type (
	httpAgreement struct {
		requestBuffer []byte
		request       *Request
	}
)

var handlerFunc func(req *Request, resp *Response)

func socketRead(conn *socket.Connection, data []byte) {
	agreement := conn.Agreement.(*httpAgreement)

	if agreement.request != nil {
		agreement.requestBuffer = append(agreement.requestBuffer, data...)

		if agreement.request.BodyNum > len(agreement.requestBuffer) {
			return
		}

		agreement.request.Body = agreement.requestBuffer[:agreement.request.BodyNum]
		agreement.requestBuffer = agreement.requestBuffer[agreement.request.BodyNum:]

		resp := productionHttpResponse(agreement.request)
		handlerFunc(agreement.request, resp)
		conn.SendData(resp.change2bytes())
		agreement.request = nil
		return
	}

	var (
		temData   []byte
		appendLen int
	)

	if len(agreement.requestBuffer) < 4 {
		appendLen = len(agreement.requestBuffer)
		temData = append(agreement.requestBuffer, data...)
	} else {
		appendLen = 3
		temData = append(agreement.requestBuffer[len(agreement.requestBuffer)-3:], data...)
	}

	headerEndNum, flag := findAgreementSpilt(temData)
	if flag == false {
		conn.Agreement = append(agreement.requestBuffer, data...)
		return
	}

	headerEndNum -= appendLen
	contentBytes := append(agreement.requestBuffer, data...)
	headerBytes := contentBytes[:headerEndNum]
	bodyBytes := contentBytes[headerEndNum:]
	agreement.request = parseHttpHeader(headerBytes)

	if agreement.request == nil {
		return
	}

	if agreement.request.BodyNum <= len(bodyBytes) {
		agreement.request.Body = bodyBytes[:agreement.request.BodyNum]

		resp := productionHttpResponse(agreement.request)
		handlerFunc(agreement.request, resp)
		conn.SendData(resp.change2bytes())
		agreement.request = nil
	}

	agreement.requestBuffer = bodyBytes
}

func socketOpen(coon *socket.Connection) {
	coon.Agreement = &httpAgreement{}
}

func socketClose(coon *socket.Connection) {
}

func findAgreementSpilt(data []byte) (int, bool) {
	for i := 0; i < len(data)-3; i++ {
		if data[i] == '\r' && data[i+1] == '\n' && data[i+2] == '\r' && data[i+3] == '\n' {
			return i + 3, true
		}
	}

	return 0, false
}

func Run(addr string, handler func(req *Request, resp *Response)) {
	handlerFunc = handler

	ser, err := server.NewServer("tcp", addr, socketOpen, socketRead, socketClose)
	if err != nil {
		panic(err)
	}

	ser.Run()
}
