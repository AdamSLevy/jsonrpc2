// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

import (
	"encoding/json"
)

// Version is the valid version string for the "jsonrpc" field required in all
// JSON RPC 2.0 objects.
const Version = "2.0"

// Response represents a JSON-RPC 2.0 Response object.
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// response hides the json.Marshaler interface that Response implements.
// Response.MarshalJSON uses this type to avoid infinite recursion.
type response Response

// MarshalJSON outputs a JSON RPC Response object with the "jsonrpc" field
// populated to Version ("2.0").
func (r Response) MarshalJSON() ([]byte, error) {
	r.JSONRPC = Version
	return json.Marshal(response(r))
}

// IsValid returns true if JSONRPC is equal to the Version ("2.0") and one of
// Response.Result or Response.Error is not nil, but not both.
func (r Response) IsValid() bool {
	return r.JSONRPC == Version &&
		(r.HasResult() && !r.HasError()) ||
		(r.HasError() && !r.HasResult())
}

func (r Response) HasError() bool {
	return r.Error != nil
}
func (r Response) HasResult() bool {
	return r.Result != nil
}

// String returns a string of the JSON with "<-- " prefixed to represent a
// Response object.
func (r Response) String() string {
	b, _ := json.Marshal(r)
	return "<-- " + string(b)
}

// BatchResponse is a type that implements fmt.Stringer for a slice of
// Responses.
type BatchResponse []Response

// String returns a string of the JSON array with "<-- " prefixed to represent
// a BatchResponse object.
func (br BatchResponse) String() string {
	s := "<-- [\n"
	for i, res := range br {
		s += "  " + res.String()[4:]
		if i < len(br)-1 {
			s += ","
		}
		s += "\n"
	}
	return s + "]"
}
