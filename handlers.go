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
	ID      int
	Title   string
	Content string
	Date    string
}

var indexPage string

// index page
func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	var posts []Post

	rows, err := DB.Query("select id, content, title, date from posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var content string
		var title string
		var date time.Time
		rows.Scan(&id, &content, &title, &date)
		posts = append(posts, Post{ID: id, Title: title, Content: content, Date: date.Format("2006-01-02")})
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
	var id = -1
	var title string
	var content string
	var date time.Time

	v := request.URL.Query()
	pID := v.Get("id")
	if len(pID) > 0 {
		stmt, err := DB.Prepare("select id, title, content, date from posts where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(pID).Scan(&id, &title, &content, &date)
		if err != nil {
			log.Fatal(err)
		}

		type Page struct {
			ID      int
			Title   string
			Content string
		}
		var page = Page{ID: id, Title: title, Content: content}

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
	var date time.Time

	v := request.URL.Query()
	pID := v.Get("id")
	if len(pID) > 0 {
		stmt, err := DB.Prepare("select id, title, content, date from posts where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(pID).Scan(&id, &title, &content, &date)
		if err != nil {
			log.Fatal(err)
		}
	}

	type Page struct {
		ID      int
		Title   string
		Content string
	}
	var page = Page{ID: id, Title: title, Content: content}

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
	pContent := request.FormValue("content")
	pDate := time.Now()

	if pID == "-1" {
		tx, err := DB.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into posts (title, content, date) values (?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		r, err := stmt.Exec(pTitle, pContent, pDate)
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
		stmt, err := tx.Prepare("update posts set title=?, content=?, date=? where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(pTitle, pContent, pDate, pID)
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
