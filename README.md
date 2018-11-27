# jsonrpc2/v7
[![GoDoc](https://godoc.org/github.com/AdamSLevy/jsonrpc2?status.svg)](https://godoc.org/github.com/AdamSLevy/jsonrpc2)
[![Go Report Card](https://goreportcard.com/badge/github.com/AdamSLevy/jsonrpc2)](https://goreportcard.com/report/github.com/AdamSLevy/jsonrpc2)
[![Coverage Status](https://coveralls.io/repos/github/AdamSLevy/jsonrpc2/badge.svg?branch=master)](https://coveralls.io/github/AdamSLevy/jsonrpc2?branch=master)
[![Build Status](https://travis-ci.org/AdamSLevy/jsonrpc2.svg?branch=master)](https://travis-ci.org/AdamSLevy/jsonrpc2)

Package `jsonrpc2` is a conforming implementation of the JSON-RPC 2.0 protocol
designed to provide a minimalist API, avoid unnecessary unmarshaling and memory
allocation, and work with any http server framework that uses http.Handler.

The `HTTPRequestHandler` will recover from any `MethodFunc` panics and will
always respond with a valid JSON RPC Response, unless of course the request was
a notification.

It strives to conform to the official specification:
[https://www.jsonrpc.org](https://www.jsonrpc.org)


## Getting started
Please read the official godoc documentation for the most up to date
information.

### Client

Clients can use the `Request`, `Response`, and `Error` types with the `json`
and `http` packages to make HTTP JSON-RPC 2.0 calls and parse their responses.
```go
reqBytes, _ := json.Marshal(jsonrpc2.NewRequest("subtract", 0, []int{5, 1}))
httpResp, _ := http.Post("www.example.com", "application/json",
        bytes.NewReader(reqBytes))
respBytes, _ := ioutil.ReadAll(httpResp.Body)
response := jsonrpc2.Response{}
json.Unmarshal(respBytes, &response)
```

### Server

Servers must implement their RPC method functions to match the `MethodFunc`
type, and relate a name to the method using a `MethodMap`.
```go
var func versionMethod(p json.RawMessage) jsonrpc2.Response {
        if p != nil {
                return jsonrpc2.NewInvalidParamsErrorResponse(nil)
        }
        return jrpc.NewResponse("0.0.0")
}
var methods = MethodMap{"version", versionMethod}
```
Read the documentation for MethodFunc and MethodMap for more information.

Finally generate an http.HandlerFunc for your MethodMap and start your server.
```go
http.ListenAndServe(":8080", jsonrpc2.HTTPRequestHandler(methods))
```
