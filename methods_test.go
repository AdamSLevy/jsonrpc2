package jsonrpc2

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodMap(t *testing.T) {
	assert := assert.New(t)

	var methods MethodMap
	assert.EqualError(methods.IsValid(), "nil MethodMap")

	methods = MethodMap{}
	assert.EqualError(methods.IsValid(), "empty MethodMap")

	methods = MethodMap{"": func(params json.RawMessage) Response { return Response{} }}
	assert.EqualError(methods.IsValid(), "empty name")

	methods = MethodMap{"test": MethodFunc(nil)}
	assert.EqualError(methods.IsValid(),
		fmt.Sprintf("nil MethodFunc for method %#v", "test"))
}

func TestMethodFuncCall(t *testing.T) {
	assert := assert.New(t)

	var fs []MethodFunc
	fs = append(fs, func(_ json.RawMessage) Response {
		return NewErrorResponse(MethodNotFoundCode, "method not found", "test data")
	}, func(_ json.RawMessage) Response {
		return Response{}
	}, func(_ json.RawMessage) Response {
		return Response{JSONRPC: "2.0",
			Error: &Error{Message: "e", Data: map[bool]bool{true: true}}}
	}, func(_ json.RawMessage) Response {
		return Response{JSONRPC: "2.0", Result: map[bool]bool{true: true}}
	})
	for _, f := range fs {
		res := f.Call(nil)
		if assert.NotNil(res.Error) {
			assert.Equal(InternalError, *res.Error)
		}
		assert.Nil(res.Result)
	}

	var f MethodFunc = func(_ json.RawMessage) Response {
		return NewErrorResponse(100, "custom", "data")
	}
	res := f.Call(nil)
	if assert.NotNil(res.Error) {
		assert.Equal(Error{
			Code:    100,
			Message: "custom",
			Data:    json.RawMessage(`"data"`),
		}, *res.Error)
	}
	assert.Nil(res.Result)

	f = func(_ json.RawMessage) Response {
		return NewInvalidParamsErrorResponse("data")
	}
	res = f.Call(nil)
	if assert.NotNil(res.Error) {
		e := InvalidParams
		e.Data = json.RawMessage(`"data"`)
		assert.Equal(e, *res.Error)
	}
	assert.Nil(res.Result)

}
