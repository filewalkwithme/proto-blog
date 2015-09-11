package main

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
)

func TestViewPostHandlerSuccess(t *testing.T) {
	var b blog
	b.loadConfig("blog.cfg")
	b.port = "8086"
	b.authorUsername = "johndoe_1"
	b.theme = "minimal"
	go b.start()

	client := &http.Client{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	client.Jar = jar

	form := url.Values{}
	os.Setenv("blog_password_"+b.authorUsername, "123456")
	form.Set("username", b.authorUsername)
	form.Set("password", "123456")
	resp, err := client.PostForm("http://localhost:"+b.port+"/login", form)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}

	form = url.Values{}
	form.Set("id", "-1")
	form.Set("title", "TestViewPostHandler")
	form.Set("short_description", "123456")
	form.Set("src_content", "123456")
	form.Set("html_content", "123456")
	resp, err = client.PostForm("http://localhost:"+b.port+"/save", form)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}

	id := ""
	stmt, err := b.DB.Prepare("select id from posts where title = 'TestViewPostHandler'")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&id)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()

	resp, err = client.Get("http://localhost:" + b.port + "/post.html?id=" + id)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}
}

func TestViewPostHandlerFailed(t *testing.T) {
	var b blog
	b.loadConfig("blog.cfg")
	b.port = "8086"
	b.authorUsername = "johndoe_1"
	b.theme = "minimal"
	go b.start()

	client := &http.Client{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	client.Jar = jar

	form := url.Values{}
	os.Setenv("blog_password_"+b.authorUsername, "123456")
	form.Set("username", b.authorUsername)
	form.Set("password", "123456")
	resp, err := client.PostForm("http://localhost:"+b.port+"/login", form)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}

	resp, err = client.Get("http://localhost:" + b.port + "/post.html?id=-1")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 500 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}
}

func TestSaveNewPostHandler(t *testing.T) {
	var b blog
	b.loadConfig("blog.cfg")
	b.port = "8086"
	b.authorUsername = "johndoe_1"
	go b.start()

	client := &http.Client{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	client.Jar = jar

	form := url.Values{}
	os.Setenv("blog_password_"+b.authorUsername, "123456")
	form.Set("username", b.authorUsername)
	form.Set("password", "123456")
	resp, err := client.PostForm("http://localhost:"+b.port+"/login", form)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}

	before := 0
	stmt, err := b.DB.Prepare("select count(id) from posts")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&before)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()

	form = url.Values{}
	form.Set("id", "-1")
	form.Set("title", "123456")
	form.Set("short_description", "123456")
	form.Set("src_content", "123456")
	form.Set("html_content", "123456")
	resp, err = client.PostForm("http://localhost:"+b.port+"/save", form)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}

	after := 0
	stmt, err = b.DB.Prepare("select count(id) from posts")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&after)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()

	if after != before+1 {
		t.Fatalf("before count(=%v) is different than after count(=%v)-1\n", before, after)
	}

}
