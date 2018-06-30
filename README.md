# jsonrpc2 - v1.1.0
[![GoDoc](https://godoc.org/github.com/AdamSLevy/jsonrpc2?status.svg)](https://godoc.org/github.com/AdamSLevy/jsonrpc2)
[![Go Report Card](https://goreportcard.com/badge/github.com/AdamSLevy/jsonrpc2)](https://goreportcard.com/report/github.com/AdamSLevy/jsonrpc2)

Package jsonrpc2 is an easy-to-use, thin, minimalist implementation of the
JSON-RPC 2.0 protocol with a handler for HTTP servers. It avoids implementing
any HTTP helper functions and instead simply provides conforming Request and
Response Types, and an http.HandlerFunc that handles single and batch Requests,
protocol errors, and recovers panics from the application's RPC method calls.
It strives to conform to the official specification: https://www.jsonrpc.org.

## Getting started
Please read the official godoc documentation for the most up to date
information.

### Client

Clients can use the Request, Response, and Error types with the json and http
packages to make HTTP JSON-RPC 2.0 calls and parse their responses.

### Server

Servers must implement their RPC method functions to match the MethodFunc type
and then register their function with a name using RegisterMethod(name,
function). Read the documentation for RegisterMethod and MethodFunc for more
information. RemarshalJSON is a convenience function for converting the
abstract params argument into a custom concrete type.
```golang
jsonrpc2.RegisterMethod("subtract", func(params interface{}) jsonrpc2.Response {
	var p []interface{}
	var ok bool
	if p, ok = params.([]interface{}); !ok {
		return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
	}
	if len(p) != 2 {
		return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
	}
	var x [2]float64
	for i := range x {
		if x[i], ok = p[i].(float64); !ok {
			return jsonrpc2.NewErrorResponse(jsonrpc2.InvalidParams)
		}
	}
	result := x[0] - x[1]
	return jsonrpc2.NewResponse(result)
})
```
After all methods are registered set up an HTTP Server with HTTPRequestHandler
as the handler.
```golang
http.ListenAndServe(":8080", jsonrpc2.HTTPRequestHandler)
```
