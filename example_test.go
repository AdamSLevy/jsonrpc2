// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	jrpc "github.com/AdamSLevy/jsonrpc2"
)

func subtract(params interface{}) jrpc.Response {
	var p []interface{}
	var ok bool
	if p, ok = params.([]interface{}); !ok {
		return jrpc.NewErrorResponse(jrpc.InvalidParams)
	}
	if len(p) != 2 {
		return jrpc.NewErrorResponse(jrpc.InvalidParams)
	}
	var x [2]float64
	for i := range x {
		if x[i], ok = p[i].(float64); !ok {
			return jrpc.NewErrorResponse(jrpc.InvalidParams)
		}
	}
	result := x[0] - x[1]
	return jrpc.NewResponse(result)
}

func alwaysPanic(params interface{}) jrpc.Response {
	panic("PANIC")
}

func Example() {
	// Register methods.
	jrpc.RegisterMethod("subtract", subtract)
	jrpc.RegisterMethod("panic", alwaysPanic)

	// Start the server.
	go func() {
		http.ListenAndServe(":8888", jrpc.HTTPRequestHandler)
	}()

	// Make requests.
	request := jrpc.NewRequest("subtract", 0, []int{5, 1})
	fmt.Println(request)
	reqBytes, _ := json.Marshal(request)
	httpResp, _ := http.Post("http://localhost:8888/v2", "application/json",
		bytes.NewReader(reqBytes))
	respBytes, _ := ioutil.ReadAll(httpResp.Body)
	response := jrpc.Response{}
	json.Unmarshal(respBytes, &response)
	fmt.Println(response)

	fmt.Println()

	request = jrpc.NewRequest("invalid", nil, nil)
	fmt.Println(request)
	reqBytes, _ = json.Marshal(request)
	httpResp, _ = http.Post("http://localhost:8888/v2", "application/json",
		bytes.NewReader(reqBytes))
	respBytes, _ = ioutil.ReadAll(httpResp.Body)
	response = jrpc.Response{}
	json.Unmarshal(respBytes, &response)
	fmt.Println(response)

	fmt.Println()

	request = jrpc.NewRequest("panic", nil, nil)
	fmt.Println(request)
	reqBytes, _ = json.Marshal(request)
	httpResp, _ = http.Post("http://localhost:8888/v2", "application/json",
		bytes.NewReader(reqBytes))
	respBytes, _ = ioutil.ReadAll(httpResp.Body)
	response = jrpc.Response{}
	json.Unmarshal(respBytes, &response)
	fmt.Println(response)

	fmt.Println()

	request = jrpc.NewNotification("invalid", nil)
	fmt.Printf("(Notification) %v\n", request)
	reqBytes, _ = json.Marshal(request)
	httpResp, _ = http.Post("http://localhost:8888/v2", "application/json",
		bytes.NewReader(reqBytes))
	respBytes, _ = ioutil.ReadAll(httpResp.Body)
	fmt.Println("Valid notifications have no response. len(respBytes):", len(respBytes))

	// Output:
	// --> {"jsonrpc":"2.0","method":"subtract","params":[5,1],"id":0}
	// <-- {"jsonrpc":"2.0","result":4,"id":0}
	//
	// --> {"jsonrpc":"2.0","method":"invalid","id":0}
	// <-- {"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":0}
	//
	// --> {"jsonrpc":"2.0","method":"panic","id":0}
	// <-- {"jsonrpc":"2.0","error":{"code":-32603,"message":"Internal error"},"id":0}
	//
	// (Notification) --> {"jsonrpc":"2.0","method":"invalid"}
	// Valid notifications have no response. len(respBytes): 0
}
