// github.com/AdamSLevy/jsonrpc2
// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

// Error represents the "error" field in a JSON-RPC 2.0 Response object.
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Official JSON-RPC 2.0 Spec Error Codes and Messages
const (
	LowestReservedErrorCode  = -32768
	ParseErrorCode           = -32700
	InvalidRequestCode       = -32600
	MethodNotFoundCode       = -32601
	InvalidParamsCode        = -32602
	InternalErrorCode        = -32603
	HighestReservedErrorCode = -32000

	ParseErrorMessage     = "Parse error"
	InvalidRequestMessage = "Invalid Request"
	MethodNotFoundMessage = "Method not found"
	InvalidParamsMessage  = "Invalid params"
	InternalErrorMessage  = "Internal error"
)

// Official Errors
var (
	// ParseError is returned to the client if a JSON is not well formed.
	ParseError = *NewError(ParseErrorCode, ParseErrorMessage, nil)
	// InvalidRequest is returned to the client if a request does not
	// conform to JSON-RPC 2.0 spec
	InvalidRequest = *NewError(InvalidRequestCode, InvalidRequestMessage, nil)
	// MethodNotFound is returned to the client if a method is called that
	// has not been registered with RegisterMethod()
	MethodNotFound = *NewError(MethodNotFoundCode, MethodNotFoundMessage, nil)
	// InvalidParams is returned to the client if a method is called with
	// an invalid "params" object. A method's function is responsible for
	// detecting and returning this error.
	InvalidParams = *NewError(InvalidParamsCode, InvalidParamsMessage, nil)
	// InternalError is returned to the client if a method function returns
	// an invalid response object.
	InternalError = *NewError(InternalErrorCode, InternalErrorMessage, nil)
)

func NewError(code int, message string, data interface{}) *Error {
	return &Error{Code: code, Message: message, Data: data}
}
