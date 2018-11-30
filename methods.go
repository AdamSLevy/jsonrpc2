// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// DebugMethodFunc controls whether additional debug information will be
// printed to stdout in the event of an internal error when a MethodFunc is
// called. This can be helpful when users are troubleshooting their
// MethodFuncs.
var DebugMethodFunc = false
var logger = log.New(os.Stderr, "", 0)

// MethodMap associates method names with MethodFuncs and is passed to
// HTTPRequestHandler() to generate a corresponding http.HandlerFunc.
type MethodMap map[string]MethodFunc

// IsValid returns nil if methods is not nil, not empty, and contains no
// entries with either an empty name or a nil MethodFunc.
func (methods MethodMap) IsValid() error {
	if methods == nil {
		return fmt.Errorf("nil MethodMap")
	}
	if len(methods) == 0 {
		return fmt.Errorf("empty MethodMap")
	}
	for name, function := range methods {
		if len(name) == 0 {
			return fmt.Errorf("empty name")
		}
		if function == nil {
			return fmt.Errorf("nil MethodFunc for method %#v", name)
		}
	}
	return nil
}

// MethodFunc is the function signature used for RPC methods. The raw JSON of
// the params of a valid Request is passed to the MethodFunc for further
// application specific unmarshaling and validation. When a MethodFunc is
// invoked by the handler, the params json.RawMessage, if not nil, is
// guaranteed to be valid JSON representing a structured JSON type.
//
// A valid MethodFunc must return a valid Response object. If MethodFunc
// panics, or if the returned Response is not valid for whatever reason, then
// an InternalError with no Data will be returned by the http.HandlerFunc
// instead.
//
// A valid Response must have either a Result or Error populated.
//
// An Error is considered populated if the Error.Message is not empty. If Error
// is populated, any Result will be ignored and the Error will be validated.
//
// A valid Error must be either InvalidParams or must use an ErrorCode outside
// of the reserved range. If the ErrorCode is InvalidParamsCode, then the
// correct InvalidParamsMessage will be set. In this case the MethodFunc only
// needs to ensure that the Message is not empty. MethodFuncs are encouraged to
// use NewInvalidParamsErrorResponse() for these errors.
//
// If you are getting InternalErrors from your method, set DebugMethodFunc to
// true for additional debug output about the cause of the internal error.
type MethodFunc func(params json.RawMessage) Response

// call is used to safely call a method from within an http.HandlerFunc. call
// wraps the actual invocation of the method so that it can recover from panics
// and validate and sanitize the returned Response. If the method panics or
// returns an invalid Response, an InternalError response is returned.
func (method MethodFunc) call(params json.RawMessage) (res Response) {
	defer func() {
		if r := recover(); r != nil {
			if DebugMethodFunc {
				logger.Printf("Internal error: %#v", r)
				logger.Printf("Params: %v", string(params))
				logger.Printf("Response: %+v", res)
			}
			res = newErrorResponse(nil, InternalError)
		}
	}()
	res = method(params)
	if res.IsError() {
		// This as an Error Response.
		// InvalidParamsCode is the only reserved ErrorCode MethodFuncs
		// are allowed to use.
		if res.Error.Code == InvalidParamsCode {
			// Ensure the correct message is used.
			res.Message = InvalidParamsMessage
		} else if res.Error.Code.IsReserved() {
			panic("Invalid Response.Error")
		}
		// Marshal the user provided Data to catch any potential errors
		// here instead of later in the http.HandlerFunc.
		data, err := json.Marshal(res.Data)
		if err != nil {
			panic("Cannot marshal Response.Error.Data")
		}
		res.Data = json.RawMessage(data)
	} else if res.Result != nil {
		// This must be a valid Response.
		// Marshal the user provided Result to catch any potential
		// errors here instead of later in the http.HandlerFunc.
		data, err := json.Marshal(res.Result)
		if err != nil {
			panic("Cannot marshal Response.Result")
		}
		res.Result = json.RawMessage(data)
	} else {
		panic("Both Response.Result and Response.Error are empty")
	}
	return
}
