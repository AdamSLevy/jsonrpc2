// github.com/AdamSLevy/jsonrpc2 v1.1.1
// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.
package jsonrpc2_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/AdamSLevy/jsonrpc2"
)

// Use the http and json packages to send a Request object.
func ExampleRequest() {
	reqBytes, _ := json.Marshal(jsonrpc2.NewRequest("subtract", 0, []int{5, 1}))
	httpResp, _ := http.Post("http://localhost:8888", "application/json", bytes.NewReader(reqBytes))
	respBytes, _ := ioutil.ReadAll(httpResp.Body)
	response := jsonrpc2.Response{}
	json.Unmarshal(respBytes, &response)
}

// Any panic will return InternalError to the user if the call was a request
// and not a Notification.
func ExampleMethodFunc_panic() {
	var alwaysPanic jsonrpc2.MethodFunc = func(params interface{}) jsonrpc2.Response {
		panic("don't worry, jsonrpc2 will recover you and return an internal error")
	}
	jsonrpc2.RegisterMethod("panic at the disco!", alwaysPanic)
}

// If a method function expects named params object with two numbers named "A"
// and "B" it could use the following anonymous struct to remarshal its given
// params argument. Note the use of pointers to detect the presence of
// individual parameters.
func ExampleRemarshalJSON_namedParams() {
	var subtract jsonrpc2.MethodFunc = func(params interface{}) jsonrpc2.Response {
		var p struct {
			A *float64
			B *float64
		}
		if err := jsonrpc2.RemarshalJSON(&p, params); err != nil ||
			p.A == nil || p.B == nil {
			return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
		}
		return jsonrpc2.NewResponse(*p.A - *p.B)
	}
	jsonrpc2.RegisterMethod("subtract", subtract)
}

// If a method function expects a params array of a single type, it can use a
// slice of that type with RemarshalJSON.
func ExampleRemarshalJSON_paramsArraySingleType() {
	jsonrpc2.RegisterMethod("subtract",
		func(params interface{}) jsonrpc2.Response {
			var p []float64
			if err := jsonrpc2.RemarshalJSON(&p, params); err != nil ||
				len(p) != 2 {
				return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
			}
			return jsonrpc2.NewResponse(p[0] - p[1])
		})
}

// If a method expects a params array of multiple types, there is no type it
// can be directly remarshaled into other than []interface{}, from there each
// individual param will need to be checked with a safe type assertion.
func ExampleRemarshalJSON_paramsArrayMultipleTypes() {
	jsonrpc2.RegisterMethod("repeat-string",
		func(params interface{}) jsonrpc2.Response {
			// Verify this is a params array of length 2.
			var p []interface{}
			if err := jsonrpc2.RemarshalJSON(&p, params); err != nil || len(p) != 2 {
				return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
			}
			// Verify that the arguments are a string and a number.
			var s string
			var ok bool
			if s, ok = p[0].(string); !ok {
				return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
			}
			var f float64
			if f, ok = p[1].(float64); !ok {
				return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
			}
			// Repeat s n times.
			var n = int(f)
			var result string
			for i := 0; i < n; i++ {
				result += s
			}
			return jsonrpc2.NewResponse(result)
		})
}
