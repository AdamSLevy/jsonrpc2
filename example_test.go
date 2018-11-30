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

	jrpc "github.com/AdamSLevy/jsonrpc2/v8"
)

var endpoint = "http://localhost:18888"

// Functions for making requests and printing the Requests and Responses.
func post(b []byte) []byte {
	httpResp, _ := http.Post(endpoint, "", bytes.NewReader(b))
	respBytes, _ := ioutil.ReadAll(httpResp.Body)
	return respBytes
}
func postNewRequest(method string, id, params interface{}) {
	postRequest(jrpc.NewRequest(method, id, params))
}
func postRequest(request interface{}) {
	fmt.Println(request)
	reqBytes, _ := json.Marshal(request)
	respBytes := post(reqBytes)
	parseResponse(respBytes)
}
func parseResponse(respBytes []byte) {
	var response interface{}
	if len(respBytes) == 0 {
		return
	} else if string(respBytes[0]) == "[" {
		response = &jrpc.BatchResponse{}
	} else {
		response = &jrpc.Response{}
	}
	json.Unmarshal(respBytes, response)
	fmt.Println(response)
	fmt.Println()
}
func postBytes(req string) {
	fmt.Println("-->", req)
	respBytes := post([]byte(req))
	parseResponse(respBytes)
}

// The RPC methods called in the JSON-RPC 2.0 specification examples.
func subtract(params json.RawMessage) jrpc.Response {
	// Parse either a params array of numbers or named numbers params.
	var a []float64
	if err := json.Unmarshal(params, &a); err == nil {
		if len(a) != 2 {
			return jrpc.NewInvalidParamsErrorResponse(
				"Invalid number of array params")
		}
		return jrpc.NewResponse(a[0] - a[1])
	}
	var p struct {
		Subtrahend *float64
		Minuend    *float64
	}
	if err := json.Unmarshal(params, &p); err != nil ||
		p.Subtrahend == nil || p.Minuend == nil {
		return jrpc.NewInvalidParamsErrorResponse("Required fields " +
			`"subtrahend" and "minuend" must be valid numbers.`)
	}
	return jrpc.NewResponse(*p.Minuend - *p.Subtrahend)
}
func sum(params json.RawMessage) jrpc.Response {
	var p []float64
	if err := json.Unmarshal(params, &p); err != nil {
		return jrpc.NewInvalidParamsErrorResponse(nil)
	}
	sum := float64(0)
	for _, x := range p {
		sum += x
	}
	return jrpc.NewResponse(sum)
}
func notifyHello(_ json.RawMessage) jrpc.Response {
	return jrpc.NewResponse("")
}
func getData(_ json.RawMessage) jrpc.Response {
	return jrpc.NewResponse([]interface{}{"hello", 5})
}

// This example makes all of the calls from the examples in the JSON-RPC 2.0
// specification and prints them in a similar format.
func Example() {
	// Start the server.
	go func() {
		// Register RPC methods.
		methods := jrpc.MethodMap{
			"subtract":     subtract,
			"sum":          sum,
			"notify_hello": notifyHello,
			"get_data":     getData,
		}
		handler := jrpc.HTTPRequestHandler(methods)
		jrpc.DebugMethodFunc = false
		http.ListenAndServe(":18888", handler)
	}()

	// Make requests.
	fmt.Println("Syntax:")
	fmt.Println("--> data sent to Server")
	fmt.Println("<-- data sent to Client")
	fmt.Println("")

	fmt.Println("rpc call with positional parameters:")
	postNewRequest("subtract", 1, []int{42, 23})
	postNewRequest("subtract", 2, []int{23, 42})

	fmt.Println("rpc call with named parameters:")
	postNewRequest("subtract", 3, map[string]int{"subtrahend": 23, "minuend": 42})
	postNewRequest("subtract", 4, map[string]int{"minuend": 42, "subtrahend": 23})

	fmt.Println("a Notification:")
	postNewRequest("update", nil, []int{1, 2, 3, 4, 5})
	postNewRequest("foobar", nil, nil)
	fmt.Println()

	fmt.Println("rpc call of non-existent method:")
	postNewRequest("foobar", "1", nil)

	fmt.Println("rpc call with invalid JSON:")
	postBytes(`{"jsonrpc":"2.0","method":"foobar,"params":"bar","baz]`)

	fmt.Println("rpc call with invalid Request object:")
	postBytes(`{"jsonrpc":"2.0","method":1,"params":"bar"}`)

	fmt.Println("rpc call Batch, invalid JSON:")
	postBytes(
		`[
  {"jsonrpc":"2.0","method":"sum","params":[1,2,4],"id":"1"},
  {"jsonrpc":"2.0","method"
]`)

	fmt.Println("rpc call with an empty Array:")
	postBytes(`[]`)

	fmt.Println("rpc call with an invalid Batch (but not empty):")
	postBytes(`[1]`)

	fmt.Println("rpc call with invalid Batch:")
	postBytes(`[1,2,3]`)

	fmt.Println("rpc call Batch:")
	postBytes(`[
  {"jsonrpc":"2.0","method":"sum","params":[1,2,4],"id":"1"},
  {"jsonrpc":"2.0","method":"notify_hello","params":[7]},
  {"jsonrpc":"2.0","method":"subtract","params":[42,23],"id":"2"},
  {"foo":"boo"},
  {"jsonrpc":"2.0","method":"foo.get","params":{"name":"myself"},"id":"5"},
  {"jsonrpc":"2.0","method":"get_data","id":"9"}
]`)
	fmt.Println("rpc call Batch (all notifications):")
	postRequest(jrpc.BatchRequest{
		jrpc.NewRequest("notify_sum", nil, []int{1, 2, 4}),
		jrpc.NewRequest("notify_hello", nil, []int{7}),
	})
	fmt.Println("<-- //Nothing is returned for all notification batches")

	// Output:
	// Syntax:
	// --> data sent to Server
	// <-- data sent to Client
	//
	// rpc call with positional parameters:
	// --> {"jsonrpc":"2.0","method":"subtract","params":[42,23],"id":1}
	// <-- {"jsonrpc":"2.0","result":19,"id":1}
	//
	// --> {"jsonrpc":"2.0","method":"subtract","params":[23,42],"id":2}
	// <-- {"jsonrpc":"2.0","result":-19,"id":2}
	//
	// rpc call with named parameters:
	// --> {"jsonrpc":"2.0","method":"subtract","params":{"minuend":42,"subtrahend":23},"id":3}
	// <-- {"jsonrpc":"2.0","result":19,"id":3}
	//
	// --> {"jsonrpc":"2.0","method":"subtract","params":{"minuend":42,"subtrahend":23},"id":4}
	// <-- {"jsonrpc":"2.0","result":19,"id":4}
	//
	// a Notification:
	// --> {"jsonrpc":"2.0","method":"update","params":[1,2,3,4,5]}
	// --> {"jsonrpc":"2.0","method":"foobar"}
	//
	// rpc call of non-existent method:
	// --> {"jsonrpc":"2.0","method":"foobar","id":"1"}
	// <-- {"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":"1"}
	//
	// rpc call with invalid JSON:
	// --> {"jsonrpc":"2.0","method":"foobar,"params":"bar","baz]
	// <-- {"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}
	//
	// rpc call with invalid Request object:
	// --> {"jsonrpc":"2.0","method":1,"params":"bar"}
	// <-- {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}
	//
	// rpc call Batch, invalid JSON:
	// --> [
	//   {"jsonrpc":"2.0","method":"sum","params":[1,2,4],"id":"1"},
	//   {"jsonrpc":"2.0","method"
	// ]
	// <-- {"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}
	//
	// rpc call with an empty Array:
	// --> []
	// <-- {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}
	//
	// rpc call with an invalid Batch (but not empty):
	// --> [1]
	// <-- [
	//   {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}
	// ]
	//
	// rpc call with invalid Batch:
	// --> [1,2,3]
	// <-- [
	//   {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null},
	//   {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null},
	//   {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null}
	// ]
	//
	// rpc call Batch:
	// --> [
	//   {"jsonrpc":"2.0","method":"sum","params":[1,2,4],"id":"1"},
	//   {"jsonrpc":"2.0","method":"notify_hello","params":[7]},
	//   {"jsonrpc":"2.0","method":"subtract","params":[42,23],"id":"2"},
	//   {"foo":"boo"},
	//   {"jsonrpc":"2.0","method":"foo.get","params":{"name":"myself"},"id":"5"},
	//   {"jsonrpc":"2.0","method":"get_data","id":"9"}
	// ]
	// <-- [
	//   {"jsonrpc":"2.0","result":7,"id":"1"},
	//   {"jsonrpc":"2.0","result":19,"id":"2"},
	//   {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request"},"id":null},
	//   {"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found"},"id":"5"},
	//   {"jsonrpc":"2.0","result":["hello",5],"id":"9"}
	// ]
	//
	// rpc call Batch (all notifications):
	// --> [
	//   {"jsonrpc":"2.0","method":"notify_sum","params":[1,2,4]},
	//   {"jsonrpc":"2.0","method":"notify_hello","params":[7]}
	// ]
	// <-- //Nothing is returned for all notification batches
}
