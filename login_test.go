package main

import (
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestLoginHandlerSuccess(t *testing.T) {
	var b blog
	b.loadConfig("blog.cfg")
	b.port = "8082"
	b.theme = "minimal"
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

	if resp.Request.URL.Path != "/" {
		t.Fatalf("wrong path returned: %v (expected: /)\n", resp.Request.URL.Path)
	}
}

func TestLoginHandlerFail(t *testing.T) {
	var b blog
	b.loadConfig("blog.cfg")
	b.port = "8083"
	b.theme = "minimal"
	b.authorUsername = "johndoe_2"
	go b.start()

	form := url.Values{}
	os.Setenv("blog_password_"+b.authorUsername, "123456")
	form.Set("username", b.authorUsername)
	form.Set("password", "111111")
	resp, err := http.PostForm("http://localhost:"+b.port+"/login", form)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.Request.URL.Path != "/admin" {
		t.Fatalf("wrong path returned: %v (expected: /admin)\n", resp.Request.URL.Path)
	}
}

func TestLoginPageHandler(t *testing.T) {
	var b blog
	b.loadConfig("blog.cfg")
	b.theme = "minimal"
	b.port = "8084"
	go b.start()

	resp, err := http.Get("http://localhost:" + b.port + "/admin")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}
}
