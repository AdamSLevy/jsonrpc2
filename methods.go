// Copyright 2018 Adam S Levy
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package jsonrpc2

import (
	"context"
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
// HTTPRequestHandler to generate a corresponding http.HandlerFunc.
//
// Method names that begin with the word rpc followed by a period character
// (U+002E or ASCII 46) are reserved for rpc-internal methods and extensions
// and MUST NOT be used for anything else. If such a method name is detected
// this will panic. No such internal rpc methods are yet defined in this
// implementation.
type MethodMap map[string]MethodFunc

// MethodFunc is the function signature used for RPC methods.
//
// MethodFuncs are invoked by the HTTPRequestHandler when a valid Request is
// parsed. MethodFuncs do not need to concern themselves with the details of
// JSON-RPC 2.0 outside of the "params" field, as all parsing and validation is
// handled by the handler.
//
// The handler will call a MethodFunc with ctx set to http.Request.Context()
// and params set to JSON from the "params" field of the Request. If "params"
// was omitted or null, params will be nil. Otherwise, params is guaranteed to
// be valid JSON that represents a JSON Object or Array.
//
// A MethodFunc is responsible for application specific parsing of params.  A
// MethodFunc should return an ErrorInvalidParams if there is any issue parsing
// expected parameters.
//
// To return a success Response to the client, a MethodFunc may return any
// non-nil, non-error interface{} that does not cause an error when passed to
// json.Marshal, as the returned interface{} will be used as the
// Response.Result.
//
// To return an Error Response to the client, a MethodFunc may return an Error
// value or pointer. Returned Errors other than InvalidParams must use an
// Error.Code outside of the reserved range, else an Internal Error is returned
// instead.
//
// If a MethodFunc panics or returns any other error, an Internal Error is
// returned to the client.
//
// For additional debug output from a MethodFunc regarding the cause of an
// Internal Error, set DebugMethodFunc to true.
type MethodFunc func(ctx context.Context, params json.RawMessage) interface{}

// call is used to safely call a method from within an http.HandlerFunc. call
// wraps the actual invocation of the method so that it can recover from panics
// and validate and sanitize the returned Response. If the method panics or
// returns an invalid Response, an Internal Error is returned.
func (method MethodFunc) call(ctx context.Context,
	name string, params json.RawMessage) (res Response) {

	var result interface{}
	defer func() {
		if err := recover(); err != nil {
			res.Error = errorInternal(nil)
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
	result = method(ctx, params)
	if result == nil {
		// MethodFuncs may not return nil since responses must always
		// include a "result" or an "error".
		panic(fmt.Errorf("method %q returned nil", name))
	}
	var methodErr Error
	switch err := result.(type) {
	case *Error:
		methodErr = *err
	case Error:
		methodErr = err
	case error:
		// MethodFuncs should not normally return a generic error. If a
		// MethodFunc intends to return an error to the client it must
		// use the Error or *Error type.
		panic(fmt.Errorf("method %q returned an unexpected error: %w",
			name, err))
	}
	// Check if this is an Error Response.
	if !methodErr.IsZero() {
		// InvalidParamsCode is the only reserved ErrorCode MethodFuncs
		// are allowed to use.
		if methodErr.Code == ErrorCodeInvalidParams {
			// Ensure the correct message is used.
			methodErr.Message = ErrorMessageInvalidParams
		} else if methodErr.Code.IsReserved() {
			panic(fmt.Errorf(
				"method %q returned an Error with a reserved ErrorCode",
				name))
		}
		if methodErr.Data != nil {
			// MethodFuncs can potentiall return types that cannot
			// be marshaled. Catch that here.
			data, err := json.Marshal(methodErr.Data)
			if err != nil {
				panic(fmt.Errorf(
					"method %q returned an Error with invalid Data: %w",
					name, err))
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
		panic(fmt.Sprintf("method %q returned an invalid result: %w",
			name, err))
	}

	// Omit null results. Can occur if the method returned
	// json.RawMessage("null").
	if string(data) == "null" {
		panic(`MethodFunc error: Result marshalled to "null"`)
	}
	res.Result = json.RawMessage(data)
	return
}
