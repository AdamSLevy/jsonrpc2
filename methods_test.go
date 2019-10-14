// Copyright 2018 Adam S Levy
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package jsonrpc2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		Func: func(_ context.Context, _ json.RawMessage) interface{} {
			return errorMethodNotFound(nil)
		},
	}, {
		Name: "nil return",
		Func: func(_ context.Context, _ json.RawMessage) interface{} {
			return nil
		},
	}, {
		Name: "error return",
		Func: func(_ context.Context, _ json.RawMessage) interface{} {
			return fmt.Errorf("not the error your looking for")
		},
	}, {
		Name: "invalid Error.Data",
		Func: func(_ context.Context, _ json.RawMessage) interface{} {
			return Error{Message: "e", Data: map[bool]bool{true: true}}
		},
	}, {
		Name: "invalid Error.Data",
		Func: func(_ context.Context, _ json.RawMessage) interface{} {
			return &Error{Message: "e", Data: map[bool]bool{true: true}}
		},
	}, {
		Name: "invalid Error.Data",
		Func: func(_ context.Context, _ json.RawMessage) interface{} {
			return Error{Message: "e"}
		},
		Error: &Error{Message: "e"},
	}, {
		Name: "invalid Result",
		Func: func(_ context.Context, _ json.RawMessage) interface{} {
			return map[bool]bool{true: true}
		},
	},
}

func TestMethodFuncCall(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	log := log.New(&buf, "", 0) // record output
	DebugMethodFunc = true

	for _, test := range testMethods {
		res := test.Func.call(context.Background(), "", nil, log)
		if test.Error == nil {
			assert.Equal(errorInternal(nil), res.Error, test.Name)
		} else {
			assert.Equal(*test.Error, res.Error, test.Name)
		}
		assert.Nil(res.Result, test.Name)
	}
	assert.Contains(string(buf.Bytes()),
		"jsonrpc2: panic running method (jsonrpc2.MethodFunc)")

	var f MethodFunc = func(_ context.Context, _ json.RawMessage) interface{} {
		return Error{100, "custom", "data"}
	}
	res := f.call(context.Background(), "", nil, log)
	if assert.NotNil(res.Error) {
		assert.Equal(Error{
			Code:    100,
			Message: "custom",
			Data:    json.RawMessage(`"data"`),
		}, res.Error)
	}
	assert.Nil(res.Result)

	f = func(_ context.Context, _ json.RawMessage) interface{} {
		return ErrorInvalidParams("data")
	}
	res = f.call(context.Background(), "", nil, log)
	if assert.NotNil(res.Error) {
		e := ErrorInvalidParams(json.RawMessage(`"data"`))
		assert.Equal(e, res.Error)
	}
	assert.Nil(res.Result)
}
