# jsonrpc2/v9
[![GoDoc](https://godoc.org/github.com/AdamSLevy/jsonrpc2?status.svg)](https://godoc.org/github.com/AdamSLevy/jsonrpc2)
[![Go Report Card](https://goreportcard.com/badge/github.com/AdamSLevy/jsonrpc2)](https://goreportcard.com/report/github.com/AdamSLevy/jsonrpc2)
[![Coverage Status](https://coveralls.io/repos/github/AdamSLevy/jsonrpc2/badge.svg?branch=master)](https://coveralls.io/github/AdamSLevy/jsonrpc2?branch=master)
[![Build Status](https://travis-ci.org/AdamSLevy/jsonrpc2.svg?branch=master)](https://travis-ci.org/AdamSLevy/jsonrpc2)

Package `jsonrpc2` is a complete and strictly conforming implementation of the
JSON-RPC 2.0 protocol for both clients and servers.

It strives to conform to the official specification:
[https://www.jsonrpc.org](https://www.jsonrpc.org)


## Getting started
Please read the official godoc documentation for the most up to date
documentation.

### Client

Clients use the provided types, optionally along with their own custom data
types for making `Requests` and parsing `Response`s. The `Request` and
`Response` types are defined so that they can accept any valid types for
`"id"`, `"params"`, and `"result"`.

Clients can use the `Request`, `Response`, and `Error` types with the `json`
and `http` packages to make HTTP JSON-RPC 2.0 calls and parse their responses.
```go
reqBytes, _ := json.Marshal(jsonrpc2.NewRequest("subtract", 0, []int{5, 1}))
httpResp, _ := http.Post("www.example.com", "application/json",
        bytes.NewReader(reqBytes))
respBytes, _ := ioutil.ReadAll(httpResp.Body)
response := jsonrpc2.Response{Result: &MyCustomResultType{}}
json.Unmarshal(respBytes, &response)
```

### Server

Servers define their own `MethodFunc`s and associate them with a method name in
a `MethodMap`. Passing the `MethodMap` to `HTTPRequestHandler()` will return a
corresponding `http.Handler` which can be used with an `http.Server`. The
`http.Handler` handles both batch and single requests, catches all protocol
errors, and recovers from any panics or invalid return values from the user
provided `MethodFunc`. `MethodFunc`s only need to catch errors related to their
function such as `InvalidParams` or any user defined errors for the RPC method.

```go
func versionMethod(params json.RawMessage) interface{} {
	if params != nil {
		return jsonrpc2.NewInvalidParamsError("no params accepted")
	}
	return "0.0.0"
}
var methods = jsonrpc2.MethodMap{"version": versionMethod}
func StartServer() {
        http.ListenAndServe(":8080", jsonrpc2.HTTPRequestHandler(methods))
}
```
