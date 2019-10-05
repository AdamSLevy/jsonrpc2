// Copyright 2018 Adam S Levy
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package jsonrpc2_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// Specify the package name to avoid goimports from reverting this
	// import to an older version.
	jsonrpc2 "github.com/AdamSLevy/jsonrpc2/v12"
)

var endpoint = "http://localhost:18888"

// Functions for making requests and printing the Requests and Responses.
func post(b []byte) []byte {
	httpResp, _ := http.Post(endpoint, "", bytes.NewReader(b))
	respBytes, _ := ioutil.ReadAll(httpResp.Body)
	return respBytes
}
func postNewRequest(method string, id, params interface{}) {
	postRequest(jsonrpc2.Request{Method: method, ID: id, Params: params})
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
		response = &jsonrpc2.BatchResponse{}
	} else {
		response = &jsonrpc2.Response{}
	}
	if err := json.Unmarshal(respBytes, response); err != nil {
		fmt.Println(string(respBytes), err)
		return
	}
	fmt.Println(response)
	fmt.Println()
}
func postBytes(req string) {
	fmt.Println("-->", req)
	respBytes := post([]byte(req))
	parseResponse(respBytes)
}

// The RPC methods called in the JSON-RPC 2.0 specification examples.
func subtract(_ context.Context, params json.RawMessage) interface{} {
	// Parse either a params array of numbers or named numbers params.
	var a []float64
	if err := json.Unmarshal(params, &a); err == nil {
		if len(a) != 2 {
			return jsonrpc2.ErrorInvalidParams("Invalid number of array params")
		}
		return a[0] - a[1]
	}
	var p struct {
		Subtrahend *float64
		Minuend    *float64
	}
	if err := json.Unmarshal(params, &p); err != nil ||
		p.Subtrahend == nil || p.Minuend == nil {
		return jsonrpc2.ErrorInvalidParams(`Required fields "subtrahend" and ` +
			`"minuend" must be valid numbers.`)
	}
	return *p.Minuend - *p.Subtrahend
}
func sum(_ context.Context, params json.RawMessage) interface{} {
	var p []float64
	if err := json.Unmarshal(params, &p); err != nil {
		return jsonrpc2.ErrorInvalidParams(err)
	}
	sum := float64(0)
	for _, x := range p {
		sum += x
	}
	return sum
}
func notifyHello(_ context.Context, _ json.RawMessage) interface{} {
	return ""
}
func getData(_ context.Context, _ json.RawMessage) interface{} {
	return []interface{}{"hello", 5}
}

// This example makes all of the calls from the examples in the JSON-RPC 2.0
// specification and prints them in a similar format.
func Example() {
	// Start the server.
	go func() {
		// Register RPC methods.
		methods := jsonrpc2.MethodMap{
			"subtract":     subtract,
			"sum":          sum,
			"notify_hello": notifyHello,
			"get_data":     getData,
		}
		jsonrpc2.DebugMethodFunc = true
		handler := jsonrpc2.HTTPRequestHandler(methods)
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
	postRequest(jsonrpc2.BatchRequest{
		jsonrpc2.Request{Method: "notify_sum", ID: nil, Params: []int{1, 2, 4}},
		jsonrpc2.Request{Method: "notify_hello", ID: nil, Params: []int{7}},
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
	// <-- {"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found","data":"foobar"},"id":"1"}
	//
	// rpc call with invalid JSON:
	// --> {"jsonrpc":"2.0","method":"foobar,"params":"bar","baz]
	// <-- {"jsonrpc":"2.0","error":{"code":-32700,"message":"Parse error"},"id":null}
	//
	// rpc call with invalid Request object:
	// --> {"jsonrpc":"2.0","method":1,"params":"bar"}
	// <-- {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request","data":"json: cannot unmarshal number into Go struct field jRequest.method of type string"},"id":null}
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
	// <-- {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request","data":"empty batch request"},"id":null}
	//
	// rpc call with an invalid Batch (but not empty):
	// --> [1]
	// <-- [
	//   {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request","data":"json: cannot unmarshal number into Go value of type jsonrpc2.jRequest"},"id":null}
	// ]
	//
	// rpc call with invalid Batch:
	// --> [1,2,3]
	// <-- [
	//   {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request","data":"json: cannot unmarshal number into Go value of type jsonrpc2.jRequest"},"id":null},
	//   {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request","data":"json: cannot unmarshal number into Go value of type jsonrpc2.jRequest"},"id":null},
	//   {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request","data":"json: cannot unmarshal number into Go value of type jsonrpc2.jRequest"},"id":null}
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
	//   {"jsonrpc":"2.0","error":{"code":-32600,"message":"Invalid Request","data":"json: unknown field \"foo\""},"id":null},
	//   {"jsonrpc":"2.0","error":{"code":-32601,"message":"Method not found","data":"foo.get"},"id":"5"},
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
