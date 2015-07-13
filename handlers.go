package main

import (
	"io/ioutil"
	"net/http"
	"text/template"
)

var indexPage string

// index page
func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	type Page struct{}
	var page = Page{}

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
