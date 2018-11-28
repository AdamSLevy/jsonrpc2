// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package jsonrpc2 is a conforming implementation of the JSON-RPC 2.0 protocol
// designed to provide a minimalist API, avoid unnecessary unmarshaling and
// memory allocation, and work with any http server framework that uses
// http.Handler. It strives to conform very strictly to the official
// specification: https://www.jsonrpc.org.
//
// This package provides types for Requests and Responses, and a function to
// return a http.HandlerFunc that calls the MethodFuncs in a given MethodMap.
// The http.HandlerFunc will recover from any MethodFunc panics and will always
// respond with a valid JSON RPC Response, unless of course the Request did not
// have an ID, and thus was a Notification.
//
// Client
//
// Clients can use the Request, Response, and Error types with the json and
// http packages to make HTTP JSON-RPC 2.0 calls and parse their responses.
//      reqBytes, _ := json.Marshal(jsonrpc2.NewRequest("subtract", 0, []int{5, 1}))
//      httpResp, _ := http.Post("www.example.com", "application/json",
//              bytes.NewReader(reqBytes))
//      respBytes, _ := ioutil.ReadAll(httpResp.Body)
//      response := jsonrpc2.Response{Result: MyCustomResultType{}}
//      json.Unmarshal(respBytes, &response)
//
// Server
//
// Servers must implement their RPC method functions to match the MethodFunc
// type, and relate a name to the method using a MethodMap.
//      var func versionMethod(p json.RawMessage) jsonrpc2.Response {
//      	if p != nil {
//      		return jsonrpc2.NewInvalidParamsErrorResponse(nil)
//      	}
//      	return jsonrpc2.NewResponse("0.0.0")
//      }
//      var methods = jsonrpc2.MethodMap{"version", versionMethod}
// Read the documentation for MethodFunc and MethodMap for more information.
//
// Finally generate an http.HandlerFunc for your MethodMap and start your
// server.
//      http.ListenAndServe(":8080", jsonrpc2.HTTPRequestHandler(methods))
package jsonrpc2
