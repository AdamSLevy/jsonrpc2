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
// application specific unmarshaling and validation. When called by the
// handler, the params json.RawMessage, if not nil, is guaranteed to be valid
// JSON representing a structured JSON type.
//
// A valid MethodFunc must return a valid Response object. If MethodFunc
// panics, or if the returned Response is not valid for whatever reason, then
// an InternalError with no Data will be returned.
//
// A valid Response must have either a Result or Error populated.
//
// If Error is populated, the Result will be discarded and the Error will be
// validated. Valid errors will always retain their Data.
//
// A valid Error must be either InvalidParams or must use an ErrorCode outside
// of the reserved range. If the ErrorCode is InvalidParamsCode, then the
// correct InvalidParamsMessage will be set, so the MethodFunc does not need to
// ensure that the Message is populated in this case. Otherwise the Message
// must be populated and the ErrorCode must not be within the reserved
// ErrorCode range.
//
// If you are getting InternalErrors from your method, set DebugMethodFunc to
// true for additional debug output about the cause of the internal error.
type MethodFunc func(params json.RawMessage) Response

// call is used to safely call a method from within an http.HandlerFunc. call
// wraps the actual invocation of the method so that it can recover from panics
// and validate and sanitize the returned Response. If the method panics or
// returns an invalid Response, an InternalError response is returned.
//
// Valid error Responses are stripped of any Result left over by the method,
// and any user provided Data is Marshaled and replaced with the resulting
// json.RawMessage.
//
// For valid Responses, the user provided Result is Marshaled and replaced with
// the resulting json.RawMessage.
//
// If you are getting InternalErrors from your method, set DebugMethodFunc to
// true for additional debug output about the cause of the internal error.
//
// See MethodFunc for more information on writing conforming methods.
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
	if !res.IsValid() {
		panic("Invalid Response")
	}
	if res.Error != nil {
		// Discard any result that may have been saved.
		res.Result = nil
		if res.Error.Code == InvalidParamsCode {
			res.Message = InvalidParamsMessage
		} else if len(res.Error.Message) == 0 || res.Error.Code.IsReserved() {
			panic("Invalid Response.Error")
		}
		data, err := json.Marshal(res.Data)
		if err != nil {
			panic("Cannot marshal Response.Error.Data")
		}
		res.Data = json.RawMessage(data)
	} else {
		data, err := json.Marshal(res.Result)
		if err != nil {
			panic("Cannot marshal Response.Result")
		}
		res.Result = json.RawMessage(data)
	}
	return
}
