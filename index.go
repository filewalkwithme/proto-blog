package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

type indexPage struct {
	BlogTitle       string
	BlogDescription string
	AdminLogged     bool
	Posts           []postEntry
}

type postEntry struct {
	ID               int
	Title            string
	Content          string
	ShortDescription string
	Author           string
	Date             string
}

// index page
func (b *blog) indexPageHandler(response http.ResponseWriter, request *http.Request) {
	var posts []postEntry

	session, err := b.store.Get(request, b.sessionName)

	if err == nil {
		rows, err := b.DB.Query("select id, html_content, short_description, title, date from posts order by date desc")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var htmlContent string
			var shortDescription string
			var title string
			var date time.Time
			rows.Scan(&id, &htmlContent, &shortDescription, &title, &date)
			posts = append(posts, postEntry{ID: id,
				Title:            title,
				Content:          htmlContent,
				ShortDescription: shortDescription,
				Author:           b.authorName,
				Date:             date.Format("2006-01-02")})
		}
		rows.Close()

		var page = indexPage{
			BlogTitle:       b.blogTitle,
			BlogDescription: b.blogDescription,
			AdminLogged:     session.Values["admin-logged"] == true,
			Posts:           posts}

		bufIndexPage, err := ioutil.ReadFile("skins/" + b.theme + "/index.html")
		if err == nil {
			indexPage := string(bufIndexPage)

			t := template.Must(template.New("indexPage").Parse(indexPage))

			t.Execute(response, page)
		} else {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}
	} else {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}
}
