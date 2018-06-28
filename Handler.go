// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package jsonrpc2

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// HTTPRequestHandler is a convenience adapter to allow the use of
// HTTPRequestHandlerFunc as an HTTP handler.
var HTTPRequestHandler = http.HandlerFunc(HTTPRequestHandlerFunc)

// HTTPRequestHandlerFunc implements an http.HandlerFunc to handle incoming
// HTTP JSON-RPC 2.0 requests. It handles both single and batch Requests,
// detects and handles ParseError, InvalidRequest, and MethodNotFound errors,
// calls the method if the request is valid and the method name has been
// registered with RegisterMethod, and returns the results of any
// non-notification Requests.
func HTTPRequestHandlerFunc(w http.ResponseWriter, req *http.Request) {
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

	rawReqs := make([]json.RawMessage, 1)
	batch := true

	// Attempt to Unmarshal as a batch request.
	if err = json.Unmarshal(reqBytes, &rawReqs); err != nil {
		batch = false
		// Attempt to Unmarshal as a single request.
		if err = json.Unmarshal(reqBytes, &rawReqs[0]); err != nil {
			// Since we know that the JSON is valid, we should
			// always be able to Unmarshal into a json.RawMessage.
			panic(err)
		}
	}

	// Catch empty batch requests.
	if len(rawReqs) == 0 {
		respondError(w, InvalidRequest)
		return
	}

	// Process all Requests.
	responses := make([]Response, 0, len(rawReqs))
	for _, rawReq := range rawReqs {
		var req Request
		var res Response
		if err = json.Unmarshal(rawReq, &req); err != nil || !req.IsValid() {
			res = NewErrorResponse(InvalidRequest)
		} else {
			if req.IsValid() {
				// Check that the method has been registered.
				method, ok := methods[req.Method]
				if !ok {
					res = newErrorResponse(req.ID, MethodNotFound)
				} else {
					res = method.Call(req.Params)
					res.ID = req.ID
				}
				// Only send a Response if the Request had an ID.
				if req.ID == nil {
					res = Response{}
				}
			} else {
				res = NewErrorResponse(InvalidRequest)
			}
		}
		if res.IsValid() {
			responses = append(responses, res)
		}
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

func respondError(w http.ResponseWriter, e Error) {
	res := NewErrorResponse(ParseError)
	respond(w, res)
}

func respond(w http.ResponseWriter, res interface{}) {
	enc := json.NewEncoder(w)
	if err := enc.Encode(res); err != nil {
		// We should never have an error encoding our Response.
		panic(err)
	}
}
