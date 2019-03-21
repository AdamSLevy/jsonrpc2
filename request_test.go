package jsonrpc2_test

import (
	"testing"

	. "github.com/AdamSLevy/jsonrpc2/v11"
	"github.com/stretchr/testify/assert"
)

var requestTests = []struct {
	Name   string
	Req    *Request
	Params interface{}
	Data   string
	Err    string
}{{
	Name:   "Request",
	Req:    &Request{ID: 5, Method: "method", Params: &[]int{0, 1}},
	Params: &[]int{},
	Data:   `{"jsonrpc":"2.0","method":"method","params":[0,1],"id":5}`,
}, {
	Name:   "Notification",
	Req:    &Request{Method: "method", Params: &struct{ D string }{D: "hi"}},
	Params: &struct{ D string }{},
	Data:   `{"jsonrpc":"2.0","method":"method","params":{"D":"hi"}}`,
}, {
	Name: "bad JSON-RPC version",
	Data: `{"method":"method","params":{"D":"hi"}}`,
	Err:  "invalid JSON-RPC 2.0 version",
}}

func TestRequest(t *testing.T) {
	t.Run("MarshalJSON", func(t *testing.T) {
		for _, test := range requestTests {
			if test.Req == nil {
				continue
			}
			t.Run(test.Name, func(t *testing.T) {
				assert := assert.New(t)
				data, err := test.Req.MarshalJSON()
				assert.NoError(err)
				assert.Equal(test.Data, string(data))
			})
		}
	})
}
