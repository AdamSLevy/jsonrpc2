package jsonrpc2

import "fmt"

var methods map[string]MethodFunc

// MethodFunc is the type of function that can be registered as an RPC method.
// When called, it will be passed a valid params object which is either
// []interface{} or map[string]interface{}. It should return a valid Response
// object with either Response.Result or Response.Error populated. If
// Response.Error is populated, Response.Result will be removed from the
// Response before sending it to the client. Any Response.Error.Code returned
// must either use the InvalidParamsCode, or use an Error.Code outside of the
// reserved range (LowestReservedErrorCode - HighestReservedErrorCode) AND have
// a non-empty Response.Error.Message. The message SHOULD be limited to a
// concise single sentence. Any additional Error.Data may also be provided.
type MethodFunc func(params interface{}) Response

// Call is used by HTTPRequestHandlerFunc to safely call a method and sanitize
// its returned Response. Invalid responses and errors are replaced by an
// InternalError response. Error responses are stripped of any Result.
func (method MethodFunc) Call(params interface{}) Response {
	r := method(params)
	if r.Error != nil {
		data := r.Error.Data
		if r.Error.Code == InvalidParamsCode {
			// Ensure the correct Error.Message is used.
			r = NewErrorResponse(InvalidParams)
		} else if len(r.Error.Message) == 0 ||
			(LowestReservedErrorCode < r.Error.Code &&
				r.Error.Code < HighestReservedErrorCode) {
			r = NewErrorResponse(InternalError)
		}
		r.Result = nil
		r.Error.Data = data
	} else if r.Result == nil {
		r = NewErrorResponse(InternalError)
	}
	return r
}

// RegisterMethod registers a new RPC method named name that calls function.
// This function is not thread safe. All RPC methods should be registered from
// a single thread and prior to serving requests with HTTPRequestHandler. This
// will return an error if either function is nil or name has already been
// registered.
func RegisterMethod(name string, function MethodFunc) error {
	if methods == nil {
		methods = make(map[string]MethodFunc)
	}
	if function == nil {
		return fmt.Errorf("methodFunc cannot be nil")
	}
	_, ok := methods[name]
	if ok {
		return fmt.Errorf("method name %v already registered", name)
	}
	methods[name] = function

	return nil
}
