package jsonrpc2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPRequestHandler(t *testing.T) {
	assert := assert.New(t)
	var methods MethodMap
	assert.Panics(func() { HTTPRequestHandler(methods) })
}
