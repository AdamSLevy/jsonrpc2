package jsonrpc2

import "encoding/json"

// Request represents a JSON RPC 2.0 Request or Notification object.
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

// NewRequest is a convenience function that returns a new Request with the
// "jsonrpc" field already populated with the required value, "2.0". If no id
// is provided, it will be populated with a default of 0. Use NewNotification
// if you want to omit the id to form a JSON RPC 2.0 Notification object.
func NewRequest(method string, id, params interface{}) Request {
	if id == nil {
		id = 0
	}
	return Request{JSONRPC: "2.0", ID: id, Method: method, Params: params}
}

// NewNotification is a convenience function that returns a new Request with no
// ID and the "jsonrpc" field already populated with the required value, "2.0".
// When a request does not have an id, it is a JSON RPC 2.0 Notification
// object.
func NewNotification(method string, params interface{}) Request {
	return Request{JSONRPC: "2.0", Method: method, Params: params}
}

// IsValid returns true when r has a valid JSONRPC value of "2.0", a
// non-empty Method, and, if not nil, valid ID and Params types.
func (r Request) IsValid() bool {
	if r.ID != nil {
		// Validate ID type.
		switch r.ID.(type) {
		case float64:
		case string:
		default:
			return false
		}
	}
	if r.Params != nil {
		// Validate Params type.
		switch r.Params.(type) {
		case []interface{}:
		case map[string]interface{}:
		default:
			return false
		}
	}
	return r.JSONRPC == "2.0" && len(r.Method) > 0
}

// String returns a JSON string with "--> " prefixed to represent a Request
// object.
func (r Request) String() string {
	b, _ := json.Marshal(r)
	return "--> " + string(b)
}
