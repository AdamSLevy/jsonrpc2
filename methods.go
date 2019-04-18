// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
)

// DebugMethodFunc controls whether additional debug information will be
// printed to stdout in the event of an InternalError when a MethodFunc is
// called. This can be helpful when troubleshooting MethodFuncs.
var DebugMethodFunc = false
var logger = log.New(os.Stdout, "", 0)

// MethodMap associates method names with MethodFuncs and is passed to
// HTTPRequestHandler() to generate a corresponding http.HandlerFunc.
type MethodMap map[string]MethodFunc

// MethodFunc is the function signature used for RPC methods. The raw JSON of
// the params of a valid Request is passed to the MethodFunc for further
// application specific unmarshaling and validation. When a MethodFunc is
// invoked by the handler, the params json.RawMessage, if not nil, is
// guaranteed to be valid JSON representing a structured JSON type. If the
// "params" field was omitted or was null, then nil is passed to the
// MethodFunc.
//
// A valid MethodFunc must return a not-nil interface{} that will not cause an
// error when passed to json.Marshal. If the underlying type of the returned
// interface{} is Error, then an Error Response will be returned to the client.
// Any return value that is not an Error will be used as the "result" field.
//
// If the MethodFunc returns an Error, then the Error must either use the
// InvalidParamsCode, or it must use an Error.Code that is outside of the
// reserved error code range. See ErrorCode.IsReserved() for more information.
//
// An invalid MethodFunc will result in an Internal Error to be returned to the
// client without revealing any information about the error. If you are getting
// InternalErrors from your MethodFunc, set DebugMethodFunc to true for
// additional debug output about the cause of the Internal Error.
type MethodFunc func(params json.RawMessage) interface{}

// call is used to safely call a method from within an http.HandlerFunc. call
// wraps the actual invocation of the method so that it can recover from panics
// and validate and sanitize the returned Response. If the method panics or
// returns an invalid Response, an InternalError response is returned.
func (method MethodFunc) call(params json.RawMessage) (res Response) {
	var result interface{}
	defer func() {
		if err := recover(); err != nil {
			res.Error = internalError(nil)
			// Clear any Result potentially left by the MethodFunc.
			res.Result = nil
			if DebugMethodFunc {
				//res.Data = err
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				logger.Printf("jsonrpc2: panic running method %#v: %v\n%s",
					method, err, buf)
				logger.Printf("jsonrpc2: Params: %v", string(params))
				logger.Printf("jsonrpc2: Return: %#v", result)
			}
		}
	}()
	result = method(params)
	if result == nil {
		// JSON-RPC 2.0 requires that Responses include a "result".
		panic(fmt.Errorf("MethodFunc: returned nil"))
	}
	var methodErr *Error
	switch err := result.(type) {
	case *Error:
		methodErr = err
	case Error:
		methodErr = &err
	case error:
		// MethodFuncs should not normally return a generic error. If a
		// MethodFunc intends to return an error to the client it must
		// use the Error or *Error type.
		panic(fmt.Errorf("MethodFunc: %v", err))
	}
	// Check if this is an Error Response.
	if methodErr != nil {
		// InvalidParamsCode is the only reserved ErrorCode MethodFuncs
		// are allowed to use.
		if methodErr.Code == InvalidParamsCode {
			// Ensure the correct message is used.
			methodErr.Message = InvalidParamsMessage
		} else if methodErr.Code.IsReserved() {
			panic("MethodFunc error: Error.Code is reserved")
		}
		if methodErr.Data != nil {
			// MethodFuncs may return types that cannot be
			// marshaled. Catch that here.
			data, err := json.Marshal(methodErr.Data)
			if err != nil {
				panic(fmt.Sprintf("MethodFunc error: Error.Data: %v", err))
			}
			// Omit null Data. Can occur if methodErr.Data is
			// json.RawMessage("null").
			if string(data) == "null" {
				data = nil
			}
			methodErr.Data = json.RawMessage(data)
		}
		res.Error = methodErr
		return
	}
	// MethodFuncs can return types that cannot be marshaled. Catch that
	// here.
	data, err := json.Marshal(result)
	if err != nil {
		panic(fmt.Sprintf("MethodFunc error: %v", err))
	}
	// Omit null Data. Can occur if methodErr.Data is
	// json.RawMessage("null").
	if string(data) == "null" {
		panic(`MethodFunc error: Result marshalled to "null"`)
	}
	res.Result = json.RawMessage(data)
	return
}
