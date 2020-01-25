package http

import (
	"strconv"
	"strings"

	"fly/common"
)

type (
	Request struct {
		Method  string
		Url     string
		Edition string
		Header  map[string]string
		Body    []byte
		BodyNum int
	}
)

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
