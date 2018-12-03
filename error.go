// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

// Error represents a JSON-RPC 2.0 Error object, which is used in the Response
// object.
type Error struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorCode represent the int JSON RPC 2.0 error code.
type ErrorCode int

// Official JSON-RPC 2.0 Spec Error Codes and Messages
const (
	LowestReservedErrorCode  ErrorCode = -32768
	ParseErrorCode           ErrorCode = -32700
	InvalidRequestCode       ErrorCode = -32600
	MethodNotFoundCode       ErrorCode = -32601
	InvalidParamsCode        ErrorCode = -32602
	InternalErrorCode        ErrorCode = -32603
	HighestReservedErrorCode ErrorCode = -32000

	ParseErrorMessage     = "Parse error"
	InvalidRequestMessage = "Invalid Request"
	MethodNotFoundMessage = "Method not found"
	InvalidParamsMessage  = "Invalid params"
	InternalErrorMessage  = "Internal error"
)

// IsReserved returns true if c is within the reserved error code range:
// [LowestReservedErrorCode, HighestReservedErrorCode].
func (c ErrorCode) IsReserved() bool {
	return LowestReservedErrorCode <= c && c <= HighestReservedErrorCode
}

// Official JSON-RPC 2.0 Errors
var (
	// ParseError is returned to the client if a JSON is not well formed.
	ParseError = NewError(ParseErrorCode, ParseErrorMessage, nil)
	// InvalidRequest is returned to the client if a request does not
	// conform to JSON-RPC 2.0 spec
	InvalidRequest = NewError(InvalidRequestCode, InvalidRequestMessage, nil)
	// MethodNotFound is returned to the client if a method is called that
	// has not been registered with RegisterMethod()
	MethodNotFound = NewError(MethodNotFoundCode, MethodNotFoundMessage, nil)
	// InvalidParams is returned to the client if a method is called with
	// an invalid "params" object. A method's function is responsible for
	// detecting and returning this error.
	InvalidParams = NewError(InvalidParamsCode, InvalidParamsMessage, nil)
	// InternalError is returned to the client if a method function returns
	// an invalid response object.
	InternalError = NewError(InternalErrorCode, InternalErrorMessage, nil)
)

// NewError returns an Error with the given code, message, and data.
func NewError(code ErrorCode, message string, data interface{}) Error {
	return Error{Code: code, Message: message, Data: data}
}

// NewInvalidParamsError returns an InvalidParams Error with the given data.
func NewInvalidParamsError(data interface{}) Error {
	err := InvalidParams
	err.Data = data
	return err
}