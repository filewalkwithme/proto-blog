package main

import (
	"net/http"
	"testing"
)

func TestIndexPageHandler(t *testing.T) {
	var b blog
	b.loadConfig("blog.cfg")
	b.theme = "minimal"
	b.port = "8081"
	go b.start()

	resp, err := http.Get("http://localhost:" + b.port)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}
}
