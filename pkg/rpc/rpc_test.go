package rpc_test

import (
	"ahmedash95/php-lsp-server/pkg/rpc"
	"testing"
)

type EncodingExample struct {
	Testing bool
}

func TestEncodeMessage(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing: true})

	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestDecodeMessage(t *testing.T) {
	message := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"
	method, content, err := rpc.DecodeMessage([]byte(message))
	contentlen := len(content)

	if err != nil {
		t.Errorf("Error decoding message: %s", err)
	}

	if method != "hi" {
		t.Errorf("Expected method test, got %s", method)
	}

	if contentlen != 15 {
		t.Errorf("Expected content length 15, got %d", contentlen)
	}
}
