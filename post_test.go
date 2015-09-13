package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestPostHandler(t *testing.T) {
	var b blog
	b.loadConfig("blog.cfg")
	b.port = "8087"
	b.authorUsername = "johndoe_1"
	b.databaseFilename = "TestPostHandler.db"
	b.theme = "minimal"
	go b.start()

	client := &http.Client{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	client.Jar = jar

	//delete all posts
	_, err = b.DB.Exec("delete from posts")
	if err != nil {
		log.Fatal(err)
	}

	//login
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

	//insert two
	form = url.Values{}
	form.Set("id", "-1")
	form.Set("title", "123456")
	form.Set("short_description", "POST 01")
	form.Set("src_content", "123456")
	form.Set("html_content", "123456")
	resp, err = client.PostForm("http://localhost:"+b.port+"/save", form)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}

	form = url.Values{}
	form.Set("id", "-1")
	form.Set("title", "123456")
	form.Set("short_description", "POST 02")
	form.Set("src_content", "123456")
	form.Set("html_content", "123456")
	resp, err = client.PostForm("http://localhost:"+b.port+"/save", form)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}

	count := 0
	stmt, err := b.DB.Prepare("select count(id) from posts")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()

	if count != 2 {
		t.Fatalf("wrong number of posts: %v (expected: 2)\n", count)
	}

	//get the max post id
	maxID := 0
	stmt, err = b.DB.Prepare("select max(id) from posts")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&maxID)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()

	//update the last one
	form = url.Values{}
	form.Set("id", strconv.Itoa(maxID))
	form.Set("title", "123456")
	form.Set("short_description", "POST UPDATED")
	form.Set("src_content", "123456")
	form.Set("html_content", "123456")
	resp, err = client.PostForm("http://localhost:"+b.port+"/save", form)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("wrong StatusCode: %v (expected: 200)\n", resp.StatusCode)
	}

	//check if the last post got updated
	resp, err = client.Get("http://localhost:" + b.port + "/post.html?id=" + strconv.Itoa(maxID))
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

	if !strings.Contains(string(body), "POST UPDATED") {
		t.Fatalf("Post not updated")
	}

	//delete the last one
	resp, err = client.Get("http://localhost:" + b.port + "/delete?id=" + strconv.Itoa(maxID))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	//count again
	count = 0
	stmt, err = b.DB.Prepare("select count(id) from posts")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer stmt.Close()

	if count != 1 {
		t.Fatalf("wrong number of posts: %v (expected: 1)\n", count)
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
