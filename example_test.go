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

	"github.com/AdamSLevy/jsonrpc2"
)

func subtract(params interface{}) jsonrpc2.Response {
	var p []interface{}
	var ok bool
	if p, ok = params.([]interface{}); !ok {
		return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
	}
	if len(p) != 2 {
		return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
	}
	var x [2]float64
	for i := range x {
		if x[i], ok = p[i].(float64); !ok {
			return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
		}
	}
	result := x[0] - x[1]
	return jsonrpc2.NewResponse(result)
}

func alwaysPanic(params interface{}) jsonrpc2.Response {
	panic("PANIC")
}

func Example() {
	// Register methods.
	jsonrpc2.RegisterMethod("subtract", subtract)
	jsonrpc2.RegisterMethod("panic", alwaysPanic)

	// Start the server.
	go func() {
		http.ListenAndServe(":8888", jsonrpc2.HTTPRequestHandler)
	}()

	// Make requests.
	request := jsonrpc2.NewRequest("subtract", 0, []int{5, 1})
	fmt.Println(request)
	reqBytes, _ := json.Marshal(request)
	httpResp, _ := http.Post("http://localhost:8888/v2", "application/json",
		bytes.NewReader(reqBytes))
	respBytes, _ := ioutil.ReadAll(httpResp.Body)
	response := jsonrpc2.Response{}
	json.Unmarshal(respBytes, &response)
	fmt.Println(response)

	fmt.Println()

	request = jsonrpc2.NewRequest("invalid", nil, nil)
	fmt.Println(request)
	reqBytes, _ = json.Marshal(request)
	httpResp, _ = http.Post("http://localhost:8888/v2", "application/json",
		bytes.NewReader(reqBytes))
	respBytes, _ = ioutil.ReadAll(httpResp.Body)
	response = jsonrpc2.Response{}
	json.Unmarshal(respBytes, &response)
	fmt.Println(response)

	fmt.Println()

	request = jsonrpc2.NewRequest("panic", nil, nil)
	fmt.Println(request)
	reqBytes, _ = json.Marshal(request)
	httpResp, _ = http.Post("http://localhost:8888/v2", "application/json",
		bytes.NewReader(reqBytes))
	respBytes, _ = ioutil.ReadAll(httpResp.Body)
	response = jsonrpc2.Response{}
	json.Unmarshal(respBytes, &response)
	fmt.Println(response)

	fmt.Println()

	request = jsonrpc2.NewNotification("invalid", nil)
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
