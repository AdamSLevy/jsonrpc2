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

// HTTPRequestHandler returns an http.HandlerFunc for the provided methods
// MethodMap. HTTPRequestHandler panics if methods.IsValid() returns an error.
//
// The returned http.HandlerFunc handles single and batch HTTP JSON-RPC 2.0
// requests. It also deals with ParseError, InvalidRequest, and MethodNotFound
// errors. For valid requests, it calls the corresponding MethodFunc and
// returns any results or errors for any non-notification Requests.
func HTTPRequestHandler(methods MethodMap) http.HandlerFunc {
	if err := methods.IsValid(); err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, req *http.Request) {
		// Read all bytes of HTTP request body.
		reqBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			respondError(w, InvalidRequest)
			return
		}

		// Check for JSON parsing issues.
		if !json.Valid(reqBytes) {
			respondError(w, ParseError)
			return
		}

		// Initially attempt to unmarshal as a batch request. Use
		// json.RawMessage at this stage so that we can parse each
		// request and return errors individually in a BatchResponse.
		rawReqs := make([]json.RawMessage, 1)
		batch := true
		if err = json.Unmarshal(reqBytes, &rawReqs); err != nil {
			// Since the JSON is known to be valid, this error
			// means it is not an array and is just a single
			// request.
			batch = false
			rawReqs[0] = json.RawMessage(reqBytes)
		}

		// Catch empty batch requests.
		if len(rawReqs) == 0 {
			respondError(w, InvalidRequest)
			return
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

		// Send any responses.
		if len(responses) > 0 {
			if batch {
				respond(w, responses)
			} else {
				respond(w, responses[0])
			}
		}
	}
}

// processRequest unmarshals and processes a single Request stored in rawReq
// using the methods defined in methods.
func processRequest(methods MethodMap, rawReq json.RawMessage) Response {
	// Use this struct to override Request.Params interface{} with a
	// json.RawMessage to avoid unnecessarily unmarshaling it as a map.
	var req struct {
		Request
		ID     json.RawMessage `json:"id"`
		Params json.RawMessage `json:"params"`
	}
	if err := unmarshalStrict(rawReq, &req); err != nil ||
		!req.IsValid() ||
		!validID(req.ID) ||
		!validParams(req.Params) {
		return newErrorResponse(json.RawMessage("null"), InvalidRequest)
	}
	// Clear null params before calling the method.
	if string(req.Params) == "null" {
		req.Params = nil
	}
	var id interface{}
	if req.ID != nil {
		id = req.ID
	}

	// Look up the requested method and call it if found.
	method, ok := methods[req.Method]
	if !ok {
		//fmt.Printf("%#v\n", req.ID)
		return newErrorResponse(id, MethodNotFound)
	}
	res := method.call(req.Params)
	res.ID = id
	return res
}

// validID assumes that id is valid JSON and returns true if id is nil, or if
// it represents a Number or a String or Null.
func validID(id json.RawMessage) bool {
	strID := string(id)
	if len(id) == 0 || strID == "null" || id[0] == '"' ||
		(id[0] != '{' && id[0] != '[' &&
			strID != "true" && strID != "false") {
		return true
	}
	return false
}

// validParams assumes that params is valid JSON and returns true if params is
// nil, or if it represents a structured value (Array or Object), or Null.
func validParams(params json.RawMessage) bool {
	if len(params) == 0 || string(params) == "null" ||
		params[0] == '[' || params[0] == '{' {
		return true
	}
	return false
}

func unmarshalStrict(data []byte, v interface{}) error {
	b := bytes.NewBuffer(data)
	d := json.NewDecoder(b)
	d.DisallowUnknownFields()
	return d.Decode(v)
}

func respondError(w http.ResponseWriter, e Error) {
	res := newErrorResponse(nil, e)
	respond(w, res)
}

func respond(w http.ResponseWriter, res interface{}) {
	enc := json.NewEncoder(w)
	// We should never have an error encoding our Response because
	// MethodFunc.call() already Marshaled the user provided Data or
	// Result. If a write error occurs there isn't anything we can do about
	// it anyway.
	enc.Encode(res)
}
