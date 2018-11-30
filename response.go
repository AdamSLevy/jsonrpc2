// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

import (
	"encoding/json"
	"fmt"
)

// Version is the valid version string for the "jsonrpc" field required in all
// JSON RPC 2.0 objects.
const Version = "2.0"

// Response represents a JSON-RPC 2.0 Response object.
//
// The json:",omitempty" are simply here for clarity. Although Error is a
// concrete type and the json package will not ever detect it as being empty,
// the json.Marhsaler interface is implemented to use the Error.Message length
// to determine whether the Error should be considered empty.
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// MarshalJSON outputs a valid JSON RPC Response object. It returns an error if
// Response.Result and Response.Error are both empty. The Error is considered
// empty if the Error.Message is empty. If the Error is not empty, any Result
// will be omitted. If the Error is empty, then the Error is omitted. The
// Response.JSONRPC field is always output as Version ("2.0").
func (r Response) MarshalJSON() ([]byte, error) {
	res := struct {
		JSONRPC string      `json:"jsonrpc"`
		Result  interface{} `json:"result,omitempty"`
		Error   interface{} `json:"error,omitempty"`
		ID      interface{} `json:"id"`
	}{JSONRPC: Version, ID: r.ID}
	if len(r.Error.Message) != 0 {
		res.Error = r.Error
		return json.Marshal(res)
	}
	if r.Result == nil {
		return nil, fmt.Errorf("Result and Error are both empty")
	}
	res.Result = r.Result
	return json.Marshal(res)
}

// NewResponse is a convenience function that returns a new success Response
// with JSONRPC already populated with the required value, "2.0".
func NewResponse(result interface{}) Response {
	return newResponse(nil, result)
}

// NewErrorResponse is a convenience function that returns a new error Response
// with JSONRPC field already populated with the required value, "2.0".
func NewErrorResponse(code ErrorCode, message string, data interface{}) Response {
	return newErrorResponse(nil, NewError(code, message, data))
}

// NewInvalidParamsErrorResponse is a convenience function that returns a
// properly formed InvalidParams error Response with the given data.
func NewInvalidParamsErrorResponse(data interface{}) Response {
	err := InvalidParams
	err.Data = data
	return newErrorResponse(nil, err)
}

func newResponse(id, result interface{}) Response {
	return Response{ID: id, Result: result}
}

func newErrorResponse(id interface{}, err Error) Response {
	return Response{ID: id, Error: err}
}

// IsValid returns true if either Response.Result is not nil or
// Response.Error.Message is not empty.
func (r Response) IsValid() bool {
	return r.JSONRPC == Version && (r.Result != nil || len(r.Error.Message) != 0)
}

// IsError returns true if r.Error.Message is not empty.
func (r Response) IsError() bool {
	return len(r.Error.Message) > 0
}

// String returns a string of the JSON with "<-- " prefixed to represent a
// Response object.
func (r Response) String() string {
	b, _ := json.Marshal(r)
	return "<-- " + string(b)
}

// BatchResponse is a type that implements String() for a slice of Responses.
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
