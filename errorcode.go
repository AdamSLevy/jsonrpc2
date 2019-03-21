// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

// ErrorCode represent the int JSON RPC 2.0 error code.
type ErrorCode int

// Official JSON-RPC 2.0 Spec Error Codes and Messages
const (
	// MinReservedErrorCode is the minimum reserved error code. Method
	// defined errors may be less than this value.
	MinReservedErrorCode ErrorCode = -32768

	// ParseErrorCode is returned to the client when invalid JSON was
	// received by the server. An error occurred on the server while
	// parsing the JSON text.
	ParseErrorCode    ErrorCode = -32700
	ParseErrorMessage           = "Parse error"

	// InvalidRequestCode is returned to the client when the JSON sent is
	// not a valid Request object.
	InvalidRequestCode    ErrorCode = -32600
	InvalidRequestMessage           = "Invalid Request"

	// MethodNotFoundCode is returned to the client when the method does
	// not exist / is not available.
	MethodNotFoundCode    ErrorCode = -32601
	MethodNotFoundMessage           = "Method not found"

	// InvalidParamsCode is returned to the client if a method is called
	// with invalid method parameter(s). MethodFuncs are responsible for
	// detecting and returning this error.
	InvalidParamsCode    ErrorCode = -32602
	InvalidParamsMessage           = "Invalid params"

	// InternalErrorCode is returned to the client if an internal error
	// occurs such as a MethodFunc panic.
	InternalErrorCode    ErrorCode = -32603
	InternalErrorMessage           = "Internal error"

	// MaxReservedErrorCode is the maximum reserved error code. Method
	// defined errors may be greater than this value.
	MaxReservedErrorCode ErrorCode = -32000
)

// IsReserved returns true if c is within the reserved error code range:
// [LowestReservedErrorCode, HighestReservedErrorCode].
func (c ErrorCode) IsReserved() bool {
	return MinReservedErrorCode <= c && c <= MaxReservedErrorCode
}
