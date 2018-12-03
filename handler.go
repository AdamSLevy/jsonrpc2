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
		return newErrorResponse(nil, InvalidRequest)
	}

	// Check for JSON parsing issues.
	if !json.Valid(reqBytes) {
		return newErrorResponse(nil, ParseError)
	}

	// Initially attempt to unmarshal as a slice to detect batch requests.
	// Use json.RawMessage at this stage so that we can parse each Request
	// and return errors individually.
	rawReqs := make([]json.RawMessage, 1)
	batch := true
	if err = json.Unmarshal(reqBytes, &rawReqs); err != nil {
		// Since the JSON is known to be valid, this error simply means
		// we have just a single request and not a batch request.
		batch = false
		rawReqs[0] = json.RawMessage(reqBytes)
	}

	// Catch empty batch requests.
	if len(rawReqs) == 0 {
		return newErrorResponse(nil, InvalidRequest)
	}

	// Process all Requests.
	responses := make(BatchResponse, 0, len(rawReqs))
	for _, rawReq := range rawReqs {
		res := processRequest(methods, rawReq)
		if res.ID == nil {
			continue
		}
		responses = append(responses, res)
	}

	// Send nothing if there are no responses.
	if len(responses) == 0 {
		return nil
	}
	// Return the entire slice if this was a batch request.
	if batch {
		return responses
	}
	return responses[0]
}

// processRequest unmarshals and processes a single Request stored in rawReq
// using the methods defined in methods.
func processRequest(methods MethodMap, rawReq json.RawMessage) Response {
	// Unmarshal and validate the Request.
	var req safeRequest
	if err := unmarshalStrict(rawReq, &req); err != nil ||
		!req.IsValid() {
		return newErrorResponse(nil, InvalidRequest)
	}
	// Clear null params before calling the method.
	if string(req.Params) == "null" {
		req.Params = nil
	}
	// Convert the ID json.RawMessage to an interface{} while respecting
	// nil values.
	var id interface{}
	if req.ID != nil {
		id = req.ID
	}

	// Look up the requested method and call it if found.
	method, ok := methods[*req.Method]
	if !ok {
		return newErrorResponse(id, MethodNotFound)
	}
	res := method.call(req.Params)
	res.ID = id

	// Log the method name if debugging is enabled and the method had an
	// InternalError.
	if DebugMethodFunc && res.Error != nil && res.Code == InternalErrorCode {
		logger.Printf("Method: %#v\n\n", *req.Method)
	}

	return res
}

// safeRequest is used to override Request.Params and Request.ID with a
// json.RawMessage to avoid unnecessarily unmarshaling it as a map. Also we use
// this to detect if the "method" field is missing.
type safeRequest struct {
	Request
	Method *string         `json:"method"`
	ID     json.RawMessage `json:"id"`
	Params json.RawMessage `json:"params"`
}

// IsValid returns true if JSONRPC is equal to the correct Version ("2.0"),
// Method is not nil, and ID and Params are the correct types.
func (r safeRequest) IsValid() bool {
	return r.JSONRPC == Version && r.Method != nil &&
		validID(r.ID) && validParams(r.Params)
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
