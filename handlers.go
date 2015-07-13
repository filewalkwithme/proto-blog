package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

var indexPage string

// index page
func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	type Post struct {
		ID   int
		Post string
		Date time.Time
	}
	var posts []Post

	rows, err := DB.Query("select id, post, date from post")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var post string
		var date time.Time
		rows.Scan(&id, &post, &date)
		posts = append(posts, Post{ID: id, Post: post, Date: date})
	}
	rows.Close()

	type Page struct {
		Posts []Post
	}
	var page = Page{Posts: posts}

	bufIndexPage, _ := ioutil.ReadFile("pages/index.html")
	indexPage = string(bufIndexPage)

	t := template.Must(template.New("page").Parse(indexPage))

	t.Execute(response, page)
}

// new post page
func postHandler(response http.ResponseWriter, request *http.Request) {
	type Page struct{}
	var page = Page{}

	bufIndexPage, _ := ioutil.ReadFile("pages/post.html")
	indexPage = string(bufIndexPage)

	t := template.Must(template.New("page").Parse(indexPage))

	t.Execute(response, page)
}
