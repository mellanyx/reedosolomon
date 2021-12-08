package main

import (
	_ "github.com/mellanyx/reedosolomon/functions"
	reedosolomon "github.com/mellanyx/reedosolomon/functions"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	// assert equality
	byte := []byte{65, 66, 67, 226, 130, 172}

	ar := reedosolomon.EncodeByteArray(byte, 285, 5)

	dec := reedosolomon.DecodeAndFixCorruptByteArray(ar, 285, 5)

	assert.Equal(t, byte, dec)
}