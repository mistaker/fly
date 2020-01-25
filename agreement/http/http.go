package http

import (
	"strconv"
	"strings"

	"fly/common"
	"fly/socket"
)

type (
	httpAgreement struct {
		requestBuffer []byte
		request       *Request
	}

	Request struct {
		Method  string
		Url     string
		Edition string
		Header  map[string]string
		Body    []byte
		BodyNum int
	}
)

func HttpSocketOpen(coon *socket.Connection) {
	coon.Agreement = &httpAgreement{}
}

func HttpSocketRead(conn *socket.Connection, data []byte) {
	agreement := conn.Agreement.(*httpAgreement)

	if agreement.request != nil {
		agreement.requestBuffer = append(agreement.requestBuffer, data...)

		if agreement.request.BodyNum > len(agreement.requestBuffer) {
			return
		}

		agreement.request.Body = agreement.requestBuffer[:agreement.request.BodyNum]
		agreement.requestBuffer = agreement.requestBuffer[agreement.request.BodyNum:]

		//todo user logic

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
		//todo user logic
	}

	agreement.requestBuffer = bodyBytes
}

func findAgreementSpilt(data []byte) (int, bool) {
	for i := 0; i < len(data)-3; i++ {
		if data[i] == '\r' && data[i+1] == '\n' && data[i+2] == '\r' && data[i+3] == '\n' {
			return i + 3, true
		}
	}

	return 0, false
}

func parseHttpHeader(data []byte) *Request {
	headerStr := common.BytesToStringFast(data)
	stringSlice := strings.Split(headerStr, "\r\n")

	if len(stringSlice) < 1 {
		return nil
	}

	firstInfo := strings.Split(stringSlice[0], " ")
	if len(firstInfo) < 3 {
		return nil
	}

	var (
		method, url, edition = firstInfo[0], firstInfo[1], firstInfo[2]
		headerS              = make(map[string]string)
		contentLen           = 0
	)

	for i := 1; i < len(stringSlice); i++ {
		temHead := strings.Split(stringSlice[i], ":")
		if len(temHead) < 2 {
			continue
		}
		headerS[temHead[0]] = temHead[1]
	}

	if method == "POST" {
		if lenStr, exit := headerS["Content-Length"]; exit {
			contentLen, _ = strconv.Atoi(lenStr)
		}
	}

	return &Request{
		Method:  method,
		Url:     url,
		Edition: edition,
		Header:  headerS,
		BodyNum: contentLen,
	}
}
