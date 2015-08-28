package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

//IndexPage struct represents the index page
type IndexPage struct {
	BlogTitle       string
	BlogDescription string
	AdminLogged     bool
	Posts           []Post
}

// index page
func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	var posts []Post

	session, err := store.Get(request, "blog-session")

	if err == nil {
		rows, err := DB.Query("select id, html_content, short_description, title, date from posts order by date desc")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var htmlContent string
			var shortDescription string
			var title string
			var date time.Time
			rows.Scan(&id, &htmlContent, &shortDescription, &title, &date)
			posts = append(posts, Post{ID: id, Title: title, Content: htmlContent, ShortDescription: shortDescription, Author: authorName, Date: date.Format("2006-01-02")})
		}
		rows.Close()

		var page = IndexPage{BlogTitle: blogTitle, BlogDescription: blogDescription, AdminLogged: session.Values["admin-logged"] == true, Posts: posts}

		bufIndexPage, err := ioutil.ReadFile("skins/" + theme + "/index.html")
		if err == nil {
			indexPage := string(bufIndexPage)

			t := template.Must(template.New("index-page").Parse(indexPage))

			t.Execute(response, page)
		} else {
			log.Printf("%v \n", err)
		}
	} else {
		log.Printf("%v \n", err)
	}
}
