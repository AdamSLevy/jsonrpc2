// github.com/AdamSLevy/jsonrpc2 v1.1.0
// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

import (
	"encoding/json"
	"fmt"
)

var methods map[string]MethodFunc

// RegisterMethod registers a new RPC method named name that calls function.
// RegisterMethod is not thread safe. All RPC methods should be registered from
// a single thread and prior to serving requests with HTTPRequestHandler. This
// will return an error if either function is nil or name has already been
// registered.
//
// See MethodFunc for more information on writing conforming methods.
func RegisterMethod(name string, function MethodFunc) error {
	if methods == nil {
		methods = make(map[string]MethodFunc)
	}
	if function == nil {
		return fmt.Errorf("methodFunc cannot be nil")
	}
	if len(name) == 0 {
		return fmt.Errorf("method name cannot be empty")
	}
	_, ok := methods[name]
	if ok {
		return fmt.Errorf("method name %v already registered", name)
	}
	methods[name] = function

	return nil
}

// MethodFunc is the type of function that can be registered as an RPC method.
// When called it will be passed a params object of either type []interface{}
// or map[string]interface{}. It should return a valid Response object with
// either Response.Result or Response.Error populated.
//
// If Response.Error is populated, Response.Result will be removed from the
// Response before sending it to the client. Any Response.Error.Code returned
// must either use the InvalidParamsCode, OR use an Error.Code outside of the
// reserved range (LowestReservedErrorCode - HighestReservedErrorCode) AND have
// a non-empty Response.Error.Message, which SHOULD be limited to a concise
// single sentence. Any additional Error.Data may also be provided.
//
// If a MethodFunc panics when it is called, or if it returns an invalid
// response, an InternalError will be sent to the client if it was not a
// Notification Request.
type MethodFunc func(params interface{}) Response

// Call is used by HTTPRequestHandlerFunc to safely call a method, recover from
// panics, and sanitize its returned Response. If method panics or returns an
// invalid response, an InternalError response is returned. Error responses are
// stripped of any Result.
//
// See MethodFunc for more information on writing conforming methods.
func (method MethodFunc) Call(params interface{}) (res Response) {
	defer func() {
		if r := recover(); r != nil {
			res = NewErrorResponse(InternalError)
		}
	}()
	res = method(params)
	if res.Error != nil {
		data := res.Error.Data
		if res.Error.Code == InvalidParamsCode {
			// Ensure the correct Error.Message is used.
			res = NewErrorResponse(InvalidParams)
		} else if len(res.Error.Message) == 0 ||
			(LowestReservedErrorCode < res.Error.Code &&
				res.Error.Code < HighestReservedErrorCode) {
			res = NewErrorResponse(InternalError)
		}
		res.Result = nil
		res.Error.Data = data
	} else if res.Result == nil {
		res = NewErrorResponse(InternalError)
	}
	return
}

// RemarshalJSON unmarshals src and remarshals it into dst as an easy way for a
// MethodFunc to marshal its provided params object into a more specific custom
// type. For this function to have any effect dst must be a pointer to a type
// that supports json.Unmarshal.
func RemarshalJSON(dst, src interface{}) error {
	var jsonBytes []byte
	var err error
	if jsonBytes, err = json.Marshal(src); err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, dst)
}
