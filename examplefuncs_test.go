// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.
package jsonrpc2_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/AdamSLevy/jsonrpc2/v10"
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
	var _ jsonrpc2.MethodFunc = func(params json.RawMessage) interface{} {
		panic("don't worry, jsonrpc2 will recover you and return an internal error")
	}
}
