package main

import (
	"testing"
)

func TestIP(t *testing.T) {
	fetchIP(hostIP{domain: "LocalHost", ip: "127.0.0.1"})
}
