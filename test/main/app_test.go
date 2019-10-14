package main

import "testing"

func TestApp(t *testing.T) {
	TestHttpGin(t)
	TestRpcGrpc(t)
}
