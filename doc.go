// Copyright 2018 Adam S Levy. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

// Package jsonrpc2 is a complete and strictly conforming implementation of the
// JSON-RPC 2.0 protocol for both clients and servers.
//
// https://www.jsonrpc.org.
//
// Client
//
// Clients use the provided types, optionally along with their own custom data
// types for making Requests and parsing Responses. The Request and Response
// types are defined so that they can accept any valid types for "id",
// "params", and "result".
//
// Clients can use the Request, Response, and Error types with the json and
// http packages to make HTTP JSON-RPC 2.0 calls and parse their responses.
//      reqBytes, _ := json.Marshal(jsonrpc2.NewRequest("subtract", 0, []int{5, 1}))
//      httpResp, _ := http.Post("www.example.com", "application/json",
//              bytes.NewReader(reqBytes))
//      respBytes, _ := ioutil.ReadAll(httpResp.Body)
//      response := jsonrpc2.Response{Result: &MyCustomResultType{}}
//      json.Unmarshal(respBytes, &response)
//
// Server
//
// Servers define their own MethodFuncs and associate them with a method name
// in a MethodMap. Passing the MethodMap to HTTPRequestHandler() will return a
// corresponding http.Handler which can be used with an http.Server. The
// http.Handler handles both batch and single requests, catches all protocol
// errors, and recovers from any panics or invalid return values from the user
// provided MethodFunc. MethodFuncs only need to catch errors related to their
// function such as Invalid Params or any user defined errors for the RPC
// method.
//
//      func versionMethod(p json.RawMessage) interface{} {
//      	if p != nil {
//      		return jsonrpc2.NewInvalidParamsError("no params accepted")
//      	}
//      	return "0.0.0"
//      }
//      var methods = jsonrpc2.MethodMap{"version": versionMethod}
//      func StartServer() {
//              http.ListenAndServe(":8080", jsonrpc2.HTTPRequestHandler(methods))
//      }
package jsonrpc2
