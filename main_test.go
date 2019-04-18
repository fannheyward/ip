package main

import (
	"net"
	"testing"
)

func TestParser(t *testing.T) {
	ips := []string{
		"127.0.0.1",
		"8.8.8.8",
		"2001:4860:4860::8888",
		"192.168.1",
	}

	for _, v := range ips {
		ip := net.ParseIP(v)
		if ip == nil {
			t.Logf("%s is NOT valid", v)
		}

		if ip.String() == v {
			t.Logf("%s is valid", v)
		}
	}
}

func TestIP(t *testing.T) {
	fetchIP("127.0.0.1")
}
