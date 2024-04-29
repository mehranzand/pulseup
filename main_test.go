package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validate_args(t *testing.T) {
	args := args{}
	assert.NotNil(t, args)
}
