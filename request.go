// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

import (
	"encoding/json"
)

// Request represents a JSON-RPC 2.0 Request or Notification object.
//
// Valid Requests must use a numeric or string type for the ID, and a
// structured type such as a slice, array, map, or struct for the Params.
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

// request hides the json.Marshaler interface that Request implements.
// Request.MarshalJSON uses this type to avoid infinite recursion.
type request Request

// MarshalJSON outputs a JSON RPC Request object with the "jsonrpc" field
// populated to Version ("2.0").
func (r Request) MarshalJSON() ([]byte, error) {
	r.JSONRPC = Version
	return json.Marshal(request(r))
}

// NewRequest returns a new Request with the given method, id, and params. If
// nil id is provided, it is by definition a Notification object and will not
// receive a response.
func NewRequest(method string, id, params interface{}) Request {
	return Request{ID: id, Method: method, Params: params}
}

// String returns a JSON string with "--> " prefixed to represent a Request
// object.
func (r Request) String() string {
	b, _ := json.Marshal(r)
	return "--> " + string(b)
}

// BatchRequest is a type that implements fmt.Stringer for a slice of Requests.
type BatchRequest []Request

// String returns a string of the JSON array with "--> " prefixed to represent
// a BatchRequest object.
func (br BatchRequest) String() string {
	s := "--> [\n"
	for i, res := range br {
		s += "  " + res.String()[4:]
		if i < len(br)-1 {
			s += ","
		}
		s += "\n"
	}
	return s + "]"
}
