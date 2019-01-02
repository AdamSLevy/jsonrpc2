// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

import "fmt"

// Error represents a JSON-RPC 2.0 Error object, which is used in the Response
// object.
type Error struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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

// Error implements the error interface.
func (e Error) Error() string {
	s := fmt.Sprintf("jsonrpc2.Error{Code:%v, Message:%#v", e.Code, e.Message)
	if e.Data != nil {
		s += fmt.Sprintf(", Data:%#v}", e.Data)
	} else {
		s += "}"
	}
	return s
}
