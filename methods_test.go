package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMethods = []struct {
	Func MethodFunc
	Name string
}{
	{
		Name: "reserved error",
		Func: func(_ json.RawMessage) interface{} {
			return MethodNotFound
		},
	}, {
		Name: "nil return",
		Func: func(_ json.RawMessage) interface{} {
			return nil
		},
	}, {

		Name: "invalid Error.Data",
		Func: func(_ json.RawMessage) interface{} {
			return Error{Message: "e", Data: map[bool]bool{true: true}}
		},
	}, {
		Name: "invalid Result",
		Func: func(_ json.RawMessage) interface{} {
			return map[bool]bool{true: true}
		},
	},
}

func TestMethodFuncCall(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	logger.SetOutput(&buf) // hide output
	DebugMethodFunc = true
	defer func() {
		logger = log.New(os.Stdout, "", 0)
	}()

	for _, test := range testMethods {
		res := test.Func.call(nil)
		if assert.NotNil(res.Error, test.Name) {
			assert.Equal(InternalError, *res.Error, test.Name)
		}
		assert.Nil(res.Result, test.Name)
	}
	assert.Equal("MethodFunc error: Error.Code is reserved\nParams: \nReturn: jsonrpc2.Error{Code:-32601, Message:\"Method not found\", Data:interface {}(nil)}\nMethodFunc error: method returned nil\nParams: \nReturn: <nil>\nMethodFunc error: Error.Data: json: unsupported type: map[bool]bool\nParams: \nReturn: jsonrpc2.Error{Code:0, Message:\"e\", Data:map[bool]bool{true:true}}\nMethodFunc error: json: unsupported type: map[bool]bool\nParams: \nReturn: map[bool]bool{true:true}\n",
		string(buf.Bytes()))

	var f MethodFunc = func(_ json.RawMessage) interface{} {
		return NewError(100, "custom", "data")
	}
	res := f.call(nil)
	if assert.NotNil(res.Error) {
		assert.Equal(Error{
			Code:    100,
			Message: "custom",
			Data:    json.RawMessage(`"data"`),
		}, *res.Error)
	}
	assert.Nil(res.Result)

	f = func(_ json.RawMessage) interface{} {
		return NewInvalidParamsError("data")
	}
	res = f.call(nil)
	if assert.NotNil(res.Error) {
		e := InvalidParams
		e.Data = json.RawMessage(`"data"`)
		assert.Equal(e, *res.Error)
	}
	assert.Nil(res.Result)
}