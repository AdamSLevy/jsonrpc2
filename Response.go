package jsonrpc2

// Response represents a JSON RPC 2.0 Response object.
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// NewResponse is a convenience function that returns a new success Response
// with the "jsonrpc" field already populated with the required value, "2.0".
func NewResponse(result interface{}) Response {
	return newResponse(nil, result)
}

// NewErrorResponse is a convenience function that returns a new error Response
// with the "jsonrpc" field already populated with the required value, "2.0".
func NewErrorResponse(err Error) Response {
	return newErrorResponse(nil, err)
}

func newResponse(id, result interface{}) Response {
	return Response{JSONRPC: "2.0", ID: id, Result: result}
}

func newErrorResponse(id interface{}, err Error) Response {
	return Response{JSONRPC: "2.0", ID: id, Error: &err}
}

// IsValid returns true when r has a valid JSONRPC value of "2.0" and one of
// Result or Error is not nil.
func (r Response) IsValid() bool {
	return r.JSONRPC == "2.0" && (r.Result != nil || r.Error != nil)
}
