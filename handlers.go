package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(secret))

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
	session, _ := store.Get(request, "blog-session")
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
		BlogTitle       string
		BlogDescription string
		AdminLogged     bool
		Posts           []Post
	}
	var page = Page{BlogTitle: blogTitle, BlogDescription: blogDescription, AdminLogged: session.Values["admin-logged"] == true, Posts: posts}

	bufIndexPage, _ := ioutil.ReadFile("pages/index.html")
	indexPage = string(bufIndexPage)

	t := template.Must(template.New("page").Parse(indexPage))

	t.Execute(response, page)
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "blog-session")
	if session.Values["admin-logged"] == true {
		http.Redirect(response, request, "/", 302)
		return
	}

	request.ParseForm()
	pUsername := request.FormValue("username")
	pPassword := request.FormValue("password")

	envPassword := os.Getenv("blog_password_" + pUsername)
	if len(envPassword) > 0 {
		if authorUsername == pUsername && envPassword == pPassword {
			session, _ := store.Get(request, "blog-session")
			session.Values["admin-logged"] = true
			session.Save(request, response)
			http.Redirect(response, request, "/", 302)
			return
		}
	}

	http.Redirect(response, request, "/login.html", 302)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "blog-session")
	if session.Values["admin-logged"] != true {
		http.Redirect(response, request, "/", 302)
		return
	}

	session.Values["admin-logged"] = false
	session.Save(request, response)

	http.Redirect(response, request, "/", 302)
}

// login page
func loginPageHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "blog-session")
	if session.Values["admin-logged"] == true {
		http.Redirect(response, request, "/", 302)
		return
	}

	type Page struct {
		BlogTitle       string
		BlogDescription string
	}
	var page = Page{BlogTitle: blogTitle, BlogDescription: blogDescription}

	bufPage, _ := ioutil.ReadFile("pages/login.html")
	t := template.Must(template.New("page").Parse(string(bufPage)))

	t.Execute(response, page)
}

// new post page
func viewPostHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "blog-session")
	var id = -1
	var title string
	var shortDescription string
	var content string
	var date time.Time

	v := request.URL.Query()
	pID := v.Get("id")
	if len(pID) > 0 {
		stmt, err := DB.Prepare("select id, title, short_description, html_content, date from posts where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		err = stmt.QueryRow(pID).Scan(&id, &title, &shortDescription, &content, &date)
		if err != nil {
			log.Fatal(err)
		}

		type Page struct {
			ID               int
			BlogTitle        string
			BlogDescription  string
			AdminLogged      bool
			Title            string
			ShortDescription string
			Author           string
			Date             string
			Content          string
		}
		var page = Page{ID: id, BlogTitle: blogTitle, BlogDescription: blogDescription, AdminLogged: session.Values["admin-logged"] == true, Title: title, ShortDescription: shortDescription, Author: authorName, Date: date.Format("2006-01-02"), Content: content}

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
	session, _ := store.Get(request, "blog-session")
	if session.Values["admin-logged"] != true {
		http.Redirect(response, request, "/", 302)
		return
	}

	var id = -1
	var title = "Title"
	var content string
	var shortDescription = "Short Description"
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
		BlogDescription  string
		AdminLogged      bool
		ShortDescription string
		Title            string
		Author           string
		Date             string
		Content          string
	}
	var page = Page{ID: id, BlogTitle: blogTitle, BlogDescription: blogDescription, AdminLogged: session.Values["admin-logged"] == true, ShortDescription: shortDescription, Title: title, Author: authorName, Date: date.Format("2006-01-02"), Content: content}

	bufIndexPage, _ := ioutil.ReadFile("pages/edit.html")
	indexPage = string(bufIndexPage)

	t := template.Must(template.New("page").Parse(indexPage))

	t.Execute(response, page)
}

// new post page
func saveHandler(response http.ResponseWriter, request *http.Request) {
	session, _ := store.Get(request, "blog-session")
	if session.Values["admin-logged"] != true {
		http.Redirect(response, request, "/", 302)
		return
	}

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
	session, _ := store.Get(request, "blog-session")
	if session.Values["admin-logged"] != true {
		http.Redirect(response, request, "/", 302)
		return
	}

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
