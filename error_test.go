package jsonrpc2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCodeIsReserved(t *testing.T) {
	assert := assert.New(t)
	var c ErrorCode
	assert.False(c.IsReserved())
	c = MinReservedErrorCode
	assert.True(c.IsReserved())
	c = MaxReservedErrorCode
	assert.True(c.IsReserved())
}

func TestError(t *testing.T) {
	assert := assert.New(t)
	var e error
	var err Error
	assert.Implements(&e, err)
}
