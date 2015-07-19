package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

//Post represents the post entry
type Post struct {
	ID               int
	Title            string
	Content          string
	ShortDescription string
	Author           string
	Date             string
}

var indexPage string

// index page
func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	var posts []Post

	rows, err := DB.Query("select id, html_content, short_description, title, date from posts")
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

	type Page struct {
		BlogTitle string
		Posts     []Post
	}
	var page = Page{BlogTitle: blogTitle, Posts: posts}

	bufIndexPage, _ := ioutil.ReadFile("pages/index.html")
	indexPage = string(bufIndexPage)

	t := template.Must(template.New("page").Parse(indexPage))

	t.Execute(response, page)
}

// new post page
func viewPostHandler(response http.ResponseWriter, request *http.Request) {
	var id = -1
	var title string
	var content string
	var date time.Time

	v := request.URL.Query()
	pID := v.Get("id")
	if len(pID) > 0 {
		stmt, err := DB.Prepare("select id, title, html_content, date from posts where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(pID).Scan(&id, &title, &content, &date)
		if err != nil {
			log.Fatal(err)
		}

		type Page struct {
			ID        int
			BlogTitle string
			Title     string
			Author    string
			Date      string
			Content   string
		}
		var page = Page{ID: id, BlogTitle: blogTitle, Title: title, Author: authorName, Date: date.Format("2006-01-02"), Content: content}

		bufIndexPage, _ := ioutil.ReadFile("pages/post.html")
		indexPage = string(bufIndexPage)

		t := template.Must(template.New("page").Parse(indexPage))

		t.Execute(response, page)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// new post page
func editHandler(response http.ResponseWriter, request *http.Request) {
	var id = -1
	var title string
	var content string
	var shortDescription string
	var date = time.Now()

	v := request.URL.Query()
	pID := v.Get("id")
	if len(pID) > 0 {
		stmt, err := DB.Prepare("select id, title, short_description, src_content, date from posts where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(pID).Scan(&id, &title, &shortDescription, &content, &date)
		if err != nil {
			log.Fatal(err)
		}
	}

	type Page struct {
		ID               int
		BlogTitle        string
		ShortDescription string
		Title            string
		Author           string
		Date             string
		Content          string
	}
	var page = Page{ID: id, BlogTitle: blogTitle, ShortDescription: shortDescription, Title: title, Author: authorName, Date: date.Format("2006-01-02"), Content: content}

	bufIndexPage, _ := ioutil.ReadFile("pages/edit.html")
	indexPage = string(bufIndexPage)

	t := template.Must(template.New("page").Parse(indexPage))

	t.Execute(response, page)
}

// new post page
func saveHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	pID := request.FormValue("id")
	pTitle := request.FormValue("title")
	pShortDescription := request.FormValue("short_description")
	pSrcContent := request.FormValue("src_content")
	pHTMLContent := request.FormValue("html_content")
	pDate := time.Now()

	if pID == "-1" {
		tx, err := DB.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into posts (title, src_content, html_content, short_description, date) values (?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		r, err := stmt.Exec(pTitle, pSrcContent, pHTMLContent, pShortDescription, pDate)
		lastID, _ := r.LastInsertId()
		pID = strconv.Itoa(int(lastID))
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
	} else {
		tx, err := DB.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("update posts set title=?, src_content=?, html_content=?, short_description=?, date=? where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(pTitle, pSrcContent, pHTMLContent, pShortDescription, pDate, pID)
		if err != nil {
			log.Fatal(err)
		}
		tx.Commit()
	}
	http.Redirect(response, request, "/edit.html?id="+pID, 302)
}

// new post page
func deleteHandler(response http.ResponseWriter, request *http.Request) {
	v := request.URL.Query()
	pID := v.Get("id")

	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("delete from posts where id=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(pID)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	http.Redirect(response, request, "/", 302)

}
