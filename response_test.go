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
	"testing"

	"github.com/stretchr/testify/assert"
)

var responseTests = []struct {
	Name   string
	Res    *Response
	Result interface{}
	Data   string
	Err    string
}{{
	Name:   "Response",
	Res:    &Response{ID: 5, Result: "result"},
	Result: "result",
	Data:   `{"jsonrpc":"2.0","result":"result","id":5}`,
}, {
	Name: "bad version",
	Data: `{"result":"result"}`,
	Err:  "invalid JSON-RPC 2.0 version",
}}

func TestResponse(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		for _, test := range responseTests {
			if test.Res == nil {
				continue
			}
			t.Run(test.Name, func(t *testing.T) {
				assert := assert.New(t)
				data, err := test.Res.MarshalJSON()
				assert.NoError(err)
				assert.Equal(test.Data, string(data))
			})
		}
	})
}
