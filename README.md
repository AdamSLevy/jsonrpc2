# jsonrpc2 - v1.0.0
[![GoDoc](https://godoc.org/github.com/AdamSLevy/jsonrpc2?status.svg)](https://godoc.org/github.com/AdamSLevy/jsonrpc2)
[![Go Report Card](https://goreportcard.com/badge/github.com/AdamSLevy/jsonrpc2)](https://goreportcard.com/report/github.com/AdamSLevy/jsonrpc2)

Package jsonrpc2 is a lightweight implementation of the JSON RPC 2.0 protocol
for HTTP clients and servers. It conforms to the [official
specification](https://www.jsonrpc.org).

## Getting started
Please read the official godoc documentation for the most up to date
information.
### Client
Clients can use the Request, Response, and Error types with the json and http
packages to make JSON RPC 2.0 calls and parse their responses.
```golang
reqBytes, _ := json.Marshal(jsonrpc2.NewRequest("subtract", 0, []int{5, 1}))
httpResp, _ := http.Post("http://localhost:8888", "application/json", bytes.NewReader(reqBytes))
respBytes, _ := ioutil.ReadAll(httpResp.Body)
response := jsonrpc2.Response{}
json.Unmarshal(respBytes, &response)
```

### Server
Servers must implement their RPC method functions to match the MethodFunc type
and then register their method with a name using RegisterMethod(name,
function). Read the documentation for RegisterMethod and MethodFunc for more
information.
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
http.ListenAndServe(":8888", jsonrpc2.HTTPRequestHandler)
```
