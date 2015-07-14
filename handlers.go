package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

type Post struct {
	ID   int
	Post string
	Date time.Time
}

var indexPage string

// index page
func indexPageHandler(response http.ResponseWriter, request *http.Request) {
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
	v := request.URL.Query()
	pID := v.Get("id")

	stmt, err := DB.Prepare("select id, post, date from post where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var id int
	var post string
	var date time.Time
	err = stmt.QueryRow(pID).Scan(&id, &post, &date)
	if err != nil {
		log.Fatal(err)
	}

	type Page struct {
		ID   int
		Post string
	}
	var page = Page{ID: id, Post: post}

	bufIndexPage, _ := ioutil.ReadFile("pages/post.html")
	indexPage = string(bufIndexPage)

	t := template.Must(template.New("page").Parse(indexPage))

	t.Execute(response, page)
}

// new post page
func saveHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	pID := request.FormValue("id")
	pPost := request.FormValue("post")
	pDate := time.Now()

	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("update post set post=?, date=? where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(pPost, pDate, pID)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()

	fmt.Printf("RequestURI: %v \n", request.Referer())
	http.Redirect(response, request, request.Referer(), 302)
}
