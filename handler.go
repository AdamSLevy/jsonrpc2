// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// HTTPRequestHandler returns an http.HandlerFunc for the given methods.
func HTTPRequestHandler(methods MethodMap) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		res := handle(methods, req)
		if res == nil {
			return
		}
		enc := json.NewEncoder(w)
		// We should never have an error encoding our Response because
		// MethodFunc.call() already Marshaled the user provided Data
		// or Result, and everything else is marshalable. If a write
		// error occurs there isn't anything we can do about it anyway.
		if err := enc.Encode(res); err != nil {
			panic(err)
		}
	}
}

// handle an http.Request for the given methods.
func handle(methods MethodMap, req *http.Request) interface{} {
	// Read all bytes of HTTP request body.
	reqBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return newErrorResponse(nil, internalError(err))
	}

	// Ensure valid JSON so it can be assumed going forward.
	if !json.Valid(reqBytes) {
		return newErrorResponse(nil, parseError(nil))
	}

	// Initially attempt to unmarshal into a slice to detect a batch
	// request. Use []json.RawMessage so that each Request can be
	// individually parsed.
	batch := true
	rawReqs := make([]json.RawMessage, 1)
	if json.Unmarshal(reqBytes, &rawReqs) != nil {
		// Since the JSON is valid, this is just a single request.
		batch = false
		rawReqs[0] = json.RawMessage(reqBytes)
	}

	// Catch empty batch requests.
	if len(rawReqs) == 0 {
		return newErrorResponse(nil, invalidRequest("empty batch request"))
	}

	// Process all Requests.
	responses := make(BatchResponse, 0, len(rawReqs))
	for _, rawReq := range rawReqs {
		res := processRequest(methods, rawReq)
		if res == nil {
			// This is a notification.
			continue
		}
		responses = append(responses, *res)
	}

	// Send nothing if there are no responses.
	if len(responses) == 0 {
		return nil
	}
	// Return a single response if this was not a batch request.
	if !batch {
		return responses[0]
	}
	return responses
}

// processRequest unmarshals and processes a single Request stored in rawReq
// using the methods defined in methods.
func processRequest(methods MethodMap, rawReq json.RawMessage) *Response {
	// Unmarshal and validate the Request.
	var id, params json.RawMessage
	req := Request{ID: &id, Params: &params}
	if err := unmarshalStrict(rawReq, &req); err != nil {
		return newErrorResponse(nil, invalidRequest(err.Error()))
	}
	if req.JSONRPC != Version {
		return newErrorResponse(nil, invalidRequest(`invalid "jsonrpc" version`))
	}
	if len(req.Method) == 0 {
		return newErrorResponse(nil, invalidRequest(`missing or empty "method"`))
	}
	if !validID(id) {
		return newErrorResponse(nil, invalidRequest(`invalid "id" type`))
	}
	if !validParams(params) {
		return newErrorResponse(nil, invalidRequest(`invalid "params" type`))
	}
	// Clear null params before calling the method.
	if string(params) == "null" {
		params = nil
	}

	// Look up the requested method and call it if found.
	method, ok := methods[req.Method]
	if !ok {
		// Don't respond to Notifications.
		if id == nil {
			return nil
		}
		return newErrorResponse(id, methodNotFound(struct {
			Method string `json:"method"`
		}{Method: req.Method}))
	}
	res := method.call(params)
	// Log the method name if debugging is enabled and the method had an
	// internal error.
	if DebugMethodFunc && res.Error != nil && res.Error.Code == InternalErrorCode {
		logger.Printf("Method: %#v\n\n", req.Method)
	}

	// Don't respond to Notifications.
	if id == nil {
		return nil
	}

	res.ID = id
	return &res
}

// validID assumes that id is valid JSON and returns true if id is nil, or if
// it represents a Number or a String or Null.
func validID(id json.RawMessage) bool {
	if len(id) == 0 || id[0] == 'n' || id[0] == '"' ||
		(id[0] != '{' && id[0] != '[' &&
			id[0] != 't' && id[0] != 'f') {
		return true
	}
	return false
}

// validParams assumes that params is valid JSON and returns true if params is
// nil, or if it represents a structured value (Array or Object), or Null.
func validParams(params json.RawMessage) bool {
	if len(params) == 0 || params[0] == 'n' ||
		params[0] == '[' || params[0] == '{' {
		return true
	}
	return false
}

// unmarshalStrict disallows unknown fields when unmarshaling JSON.
func unmarshalStrict(data []byte, v interface{}) error {
	b := bytes.NewBuffer(data)
	d := json.NewDecoder(b)
	d.DisallowUnknownFields()
	return d.Decode(v)
}

// newErrorResponse returns a Response with the ID and Error populated.
func newErrorResponse(id interface{}, err *Error) *Response {
	return &Response{ID: id, Error: err}
}
