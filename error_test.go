package jsonrpc2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorCodeIsReserved(t *testing.T) {
	assert := assert.New(t)
	var c ErrorCode
	assert.False(c.IsReserved())
	c = LowestReservedErrorCode
	assert.True(c.IsReserved())
	c = HighestReservedErrorCode
	assert.True(c.IsReserved())
}

func TestError(t *testing.T) {
	assert := assert.New(t)
	var e error
	var err Error
	assert.Implements(&e, err)
}
