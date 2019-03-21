package jsonrpc2_test

import (
	"testing"

	. "github.com/AdamSLevy/jsonrpc2/v11"
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
