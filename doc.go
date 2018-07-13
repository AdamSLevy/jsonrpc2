// github.com/AdamSLevy/jsonrpc2 v2.0.0
// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package jsonrpc2 is an easy-to-use, thin, minimalist implementation of the
// JSON-RPC 2.0 protocol with a handler for HTTP servers. It avoids
// implementing any HTTP helper functions and instead simply provides
// conforming Request and Response Types, and an http.HandlerFunc that handles
// single and batch Requests, protocol errors, and recovers panics from the
// application's RPC method calls. It strives to conform to the official
// specification: https://www.jsonrpc.org.
//
// Client
//
// Clients can use the Request, Response, and Error types with the json and
// http packages to make HTTP JSON-RPC 2.0 calls and parse their responses.
//      reqBytes, _ := json.Marshal(jsonrpc2.NewRequest("subtract", 0, []int{5, 1}))
//      httpResp, _ := http.Post("www.example.com", "application/json",
//              bytes.NewReader(reqBytes))
//      respBytes, _ := ioutil.ReadAll(httpResp.Body)
//      response := jsonrpc2.Response{}
//      json.Unmarshal(respBytes, &response)
//
// Server
//
// Servers must implement their RPC method functions to match the MethodFunc
// type.
//      type MethodFunc func(params interface{}) Response
// Methods must be registered with a name using RegisterMethod().
//      jsonrpc2.RegisterMethod("subtract", mySubtractMethodFunc)
// Read the documentation for RegisterMethod and MethodFunc for more
// information.
//
// For convenience, methods can use RemarshalJSON() for converting the abstract
// params argument into a custom concrete type.
//
// After all methods are registered, set up an HTTP Server with
// HTTPRequestHandler as the handler.
//      http.ListenAndServe(":8080", jsonrpc2.HTTPRequestHandler)
package jsonrpc2
