package main

import (
	"testing"
)

func TestParser(t *testing.T) {
	ips, err := parse("127.0.0.1")
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(ips) != 1 {
		t.Fatal("length error")
	}

	if ips[0] == "127.0.0.1" {
		t.Log("127.0.0.1 OK")
	}
}

func TestIP(t *testing.T) {
	fetchIP("127.0.0.1")
}

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parse("127.0.0.1")
	}
}
