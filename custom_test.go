package main

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCustomHandler(t *testing.T) {
	var b blog
	b.loadConfig("blog.cfg")
	b.theme = "minimal"
	b.port = "8088"
	go b.start()

	resp, err := http.Get("http://localhost:" + b.port + "/about.html")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if string(body) != "minimal theme!\n" {
		t.Fatalf("Wrong content in About.html\n")
	}
}
