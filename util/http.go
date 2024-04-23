package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

func DumpRequest(debug bool, req *http.Request) {
	if !debug {
		return
	}

	rc := *req

	data, _ := io.ReadAll(req.Body)
	rc.Body = io.NopCloser(bytes.NewBuffer(data))

	fmt.Println("---------- DEBUG :: Request")
	dd, _ := httputil.DumpRequest(&rc, true)
	fmt.Println(string(dd))

	req.Body = io.NopCloser(bytes.NewBuffer(data))
}

func DumpResponse(debug bool, resp *http.Response) {
	if !debug {
		return
	}

	rc := *resp

	data, _ := io.ReadAll(resp.Body)
	rc.Body = io.NopCloser(bytes.NewBuffer(data))

	fmt.Println("---------- DEBUG :: Response")
	dd, _ := httputil.DumpResponse(&rc, true)
	fmt.Println(string(dd))

	resp.Body = io.NopCloser(bytes.NewBuffer(data))
}
