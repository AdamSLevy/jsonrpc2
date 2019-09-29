package jsonrpc2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMethods = []struct {
	Func  MethodFunc
	Name  string
	Error *Error
}{
	{
		Name: "reserved error",
		Func: func(_ json.RawMessage) interface{} {
			return methodNotFound(nil)
		},
	}, {
		Name: "nil return",
		Func: func(_ json.RawMessage) interface{} {
			return nil
		},
	}, {
		Name: "error return",
		Func: func(_ json.RawMessage) interface{} {
			return fmt.Errorf("not the error your looking for")
		},
	}, {
		Name: "invalid Error.Data",
		Func: func(_ json.RawMessage) interface{} {
			return Error{Message: "e", Data: map[bool]bool{true: true}}
		},
	}, {
		Name: "invalid Error.Data",
		Func: func(_ json.RawMessage) interface{} {
			return &Error{Message: "e", Data: map[bool]bool{true: true}}
		},
	}, {
		Name: "invalid Error.Data",
		Func: func(_ json.RawMessage) interface{} {
			return Error{Message: "e"}
		},
		Error: &Error{Message: "e"},
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
			if test.Error == nil {
				assert.Equal(internalError(nil), res.Error, test.Name)
			} else {
				assert.Equal(test.Error, res.Error, test.Name)
			}

		}
		assert.Nil(res.Result, test.Name)
	}
	assert.Contains(string(buf.Bytes()),
		"jsonrpc2: panic running method (jsonrpc2.MethodFunc)")

	var f MethodFunc = func(_ json.RawMessage) interface{} {
		return Error{100, "custom", "data"}
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
		return InvalidParams("data")
	}
	res = f.call(nil)
	if assert.NotNil(res.Error) {
		e := InvalidParams(json.RawMessage(`"data"`))
		assert.Equal(&e, res.Error)
	}
	assert.Nil(res.Result)
}
