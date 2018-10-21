# jsonrpc2/v4
[![GoDoc](https://godoc.org/github.com/AdamSLevy/jsonrpc2?status.svg)](https://godoc.org/github.com/AdamSLevy/jsonrpc2)
[![Go Report Card](https://goreportcard.com/badge/github.com/AdamSLevy/jsonrpc2)](https://goreportcard.com/report/github.com/AdamSLevy/jsonrpc2)
[![Coverage Status](https://coveralls.io/repos/github/AdamSLevy/jsonrpc2/badge.svg?branch=master)](https://coveralls.io/github/AdamSLevy/jsonrpc2?branch=master)
[![Build Status](https://travis-ci.org/AdamSLevy/jsonrpc2.svg?branch=master)](https://travis-ci.org/AdamSLevy/jsonrpc2)

Package jsonrpc2 is a minimalist implementation of the JSON-RPC 2.0 protocol
that provides types for Requests and Responses, and an http.Handler that calls
MethodFuncs registered with RegisterMethod(). The HTTPRequestHandler will
recover from any MethodFunc panics and will always respond with a valid JSON
RPC Response, unless of course the request was a notification.

It strives to conform to the official specification:
[https://www.jsonrpc.org](https://www.jsonrpc.org)


## Getting started
Please read the official godoc documentation for the most up to date
information.

### Client

Clients can use the Request, Response, and Error types with the json and http
packages to make HTTP JSON-RPC 2.0 calls and parse their responses.
```go
reqBytes, _ := json.Marshal(jsonrpc2.NewRequest("subtract", 0, []int{5, 1}))
httpResp, _ := http.Post("www.example.com", "application/json",
        bytes.NewReader(reqBytes))
respBytes, _ := ioutil.ReadAll(httpResp.Body)
response := &jsonrpc2.Response{}
json.Unmarshal(respBytes, response)
```

### Server

Servers must implement their RPC method functions to match the MethodFunc type.
Methods must be registered with a name using RegisterMethod().
```go
     var func versionMethod(p json.RawMessage) *jsonrpc2.Response {
     	if p != nil {
     		return jsonrpc2.NewInvalidParamsErrorResponse(nil)
     	}
     	return jrpc.NewResponse("0.0.0")
     }
     jsonrpc2.RegisterMethod("version", jsonrpc2.MethodFunc(versionMethod))
```
Read the documentation for RegisterMethod and MethodFunc for more information.

After all methods are registered, set up an HTTP Server with
```go
HTTPRequestHandler as the handler.
     http.ListenAndServe(":8080", jsonrpc2.HTTPRequestHandler)
```
