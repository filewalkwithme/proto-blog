package main

import (
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestLogoutHandler(t *testing.T) {
	var b blog
	b.loadConfig("blog.cfg")
	b.port = "8085"
	b.authorUsername = "johndoe_1"
	go b.start()

	form := url.Values{}
	os.Setenv("blog_password_"+b.authorUsername, "123456")
	form.Set("username", b.authorUsername)
	form.Set("password", "123456")
	resp, err := http.PostForm("http://localhost:"+b.port+"/login", form)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}

	resp, err = http.Get("http://localhost:" + b.port + "/logout")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}
}
