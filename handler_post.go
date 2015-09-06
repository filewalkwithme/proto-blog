package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type postPage struct {
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

func (b *blog) viewPostHandler(response http.ResponseWriter, request *http.Request) {
	session, err := b.store.Get(request, b.sessionName)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	var id = -1
	var title string
	var shortDescription string
	var content string
	var date time.Time

	v := request.URL.Query()
	pID := v.Get("id")
	if len(pID) > 0 {
		stmt, err := b.DB.Prepare("select id, title, short_description, html_content, date from posts where id = ?")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}
		defer stmt.Close()
		err = stmt.QueryRow(pID).Scan(&id, &title, &shortDescription, &content, &date)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}

		var page = postPage{
			ID:               id,
			BlogTitle:        b.blogTitle,
			BlogDescription:  b.blogDescription,
			AdminLogged:      session.Values["admin-logged"] == true,
			Title:            title,
			ShortDescription: shortDescription,
			Author:           b.authorName,
			Date:             date.Format("2006-01-02"),
			Content:          content}

		bufIndexPage, err := ioutil.ReadFile("skins/" + b.theme + "/post.html")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}

		indexPage := string(bufIndexPage)

		t := template.Must(template.New("postPage").Parse(indexPage))

		t.Execute(response, page)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func (b *blog) editPostHandler(response http.ResponseWriter, request *http.Request) {
	session, err := b.store.Get(request, b.sessionName)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

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
		stmt, err := b.DB.Prepare("select id, title, short_description, src_content, date from posts where id = ? order by date desc")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}
		defer stmt.Close()
		err = stmt.QueryRow(pID).Scan(&id, &title, &shortDescription, &content, &date)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}
	}

	var page = postPage{
		ID:               id,
		BlogTitle:        b.blogTitle,
		BlogDescription:  b.blogDescription,
		AdminLogged:      session.Values["admin-logged"] == true,
		ShortDescription: shortDescription,
		Title:            title,
		Author:           b.authorName,
		Date:             date.Format("2006-01-02"),
		Content:          content}

	bufIndexPage, err := ioutil.ReadFile("skins/" + b.theme + "/edit.html")
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}
	indexPage := string(bufIndexPage)
	indexPage = strings.Replace(indexPage, "{{.Editor}}", editor, -1)

	t := template.Must(template.New("editPage").Parse(indexPage))

	t.Execute(response, page)
}

func (b *blog) savePostHandler(response http.ResponseWriter, request *http.Request) {
	session, err := b.store.Get(request, b.sessionName)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

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
		tx, err := b.DB.Begin()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}

		stmt, err := tx.Prepare("insert into posts (title, src_content, html_content, short_description, date) values (?, ?, ?, ?, ?)")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}
		defer stmt.Close()

		r, err := stmt.Exec(pTitle, pSrcContent, pHTMLContent, pShortDescription, pDate)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}

		lastID, err := r.LastInsertId()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}

		pID = strconv.Itoa(int(lastID))
		tx.Commit()
	} else {

		tx, err := b.DB.Begin()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}

		stmt, err := tx.Prepare("update posts set title=?, src_content=?, html_content=?, short_description=?, date=? where id = ?")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(pTitle, pSrcContent, pHTMLContent, pShortDescription, pDate, pID)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			log.Printf("%v \n", err)
			fmt.Fprintf(response, "%v \n", err)
			return
		}

		tx.Commit()
	}
	http.Redirect(response, request, "/edit.html?id="+pID, 302)
}

func (b *blog) deletePostHandler(response http.ResponseWriter, request *http.Request) {
	session, err := b.store.Get(request, b.sessionName)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	if session.Values["admin-logged"] != true {
		http.Redirect(response, request, "/", 302)
		return
	}

	v := request.URL.Query()
	pID := v.Get("id")

	tx, err := b.DB.Begin()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	stmt, err := tx.Prepare("delete from posts where id=?")
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(pID)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}
	tx.Commit()

	http.Redirect(response, request, "/", 302)
}
