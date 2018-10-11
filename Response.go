// github.com/AdamSLevy/jsonrpc2
// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

import "encoding/json"

// Response represents a JSON-RPC 2.0 Response object.
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// NewResponse is a convenience function that returns a new success Response
// with JSONRPC already populated with the required value, "2.0".
func NewResponse(result interface{}) *Response {
	return newResponse(nil, result)
}

// NewErrorResponse is a convenience function that returns a new error Response
// with JSONRPC field already populated with the required value, "2.0".
func NewErrorResponse(code int, message string, data interface{}) *Response {
	return newErrorResponse(nil, newError(code, message, data))
}

// NewInvalidParamsErrorResponse is a convenience function that returns a
// properly formed InvalidParams error Response with the given data.
func NewInvalidParamsErrorResponse(data interface{}) *Response {
	err := ParseError
	err.Data = data
	return newErrorResponse(nil, err)
}

func newResponse(id, result interface{}) *Response {
	return &Response{JSONRPC: "2.0", ID: id, Result: result}
}

func newErrorResponse(id interface{}, err *Error) *Response {
	return &Response{JSONRPC: "2.0", ID: id, Error: err}
}

// IsValid returns true when r has a valid JSONRPC value of "2.0" and one of
// Result or Error is not nil.
func (r Response) IsValid() bool {
	return r.JSONRPC == "2.0" && (r.Result != nil || r.Error != nil)
}

// String returns a string of the JSON with "<-- " prefixed to represent a
// Response object.
func (r Response) String() string {
	b, _ := json.Marshal(r)
	return "<-- " + string(b)
}

// BatchResponse is a type that implements String() for a slice of Responses.
type BatchResponse []*Response

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
