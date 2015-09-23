package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	ini "github.com/Unknwon/goconfig"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type blog struct {
	DB               *sql.DB
	cfg              *ini.ConfigFile
	blogTitle        string
	blogDescription  string
	authorName       string
	authorUsername   string
	secret           string
	sessionName      string
	theme            string
	port             string
	databaseFilename string

	store  *sessions.CookieStore
	router *mux.Router
}

func (b *blog) loadConfig(configFile string) {
	cfg, err := ini.LoadConfigFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	b.authorName, err = cfg.GetValue("author", "name")
	if err != nil {
		log.Fatal(err)
	}

	b.authorUsername, err = cfg.GetValue("author", "username")
	if err != nil {
		log.Fatal(err)
	}

	b.blogTitle, err = cfg.GetValue("blog", "title")
	if err != nil {
		log.Fatal(err)
	}

	b.blogDescription, err = cfg.GetValue("blog", "description")
	if err != nil {
		log.Fatal(err)
	}

	b.secret, err = cfg.GetValue("blog", "secret")
	if err != nil {
		log.Fatal(err)
	}

	if b.secret == "something-very-secret" {
		err := fmt.Errorf("Please, change the parameter 'secret' value to something different than it's default value (something-very-secret)")
		log.Fatal(err)
	}

	b.sessionName, err = cfg.GetValue("blog", "session-name")
	if err != nil {
		log.Fatal(err)
	}

	if b.sessionName == "blog-session" {
		err := fmt.Errorf("Please, change the parameter 'session-name' value to something different than it's default value (blog-session)")
		log.Fatal(err)
	}

	b.theme, err = cfg.GetValue("blog", "theme")
	if err != nil {
		log.Fatal(err)
	}

	b.port, err = cfg.GetValue("blog", "port")
	if err != nil {
		log.Fatal(err)
	}

	b.databaseFilename, err = cfg.GetValue("blog", "database-filename")
	if err != nil {
		log.Fatal(err)
	}

	_, errCheckExists := os.Stat("./" + b.databaseFilename)

	db, err := sql.Open("sqlite3", "./"+b.databaseFilename)
	if err != nil {
		log.Fatal(err)
	}
	b.DB = db

	if os.IsNotExist(errCheckExists) {
		log.Printf("File blog.db not exists. \nCreating initial database.\n")

		sqlCreateDB := `
		create table posts (id integer not null primary key, title text, src_content text, html_content text, short_description text, date datetime);
		`

		_, err = db.Exec(sqlCreateDB)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlCreateDB)
			return
		}
		log.Printf("Initial database created.\n")
	}
}

func (b *blog) start() {
	b.router = mux.NewRouter()
	b.store = sessions.NewCookieStore([]byte(b.secret))

	b.router.HandleFunc("/", b.indexPageHandler).Methods("GET")
	b.router.HandleFunc("/admin", b.loginPageHandler).Methods("GET")
	b.router.HandleFunc("/login", b.loginHandler).Methods("POST")
	b.router.HandleFunc("/logout", b.logoutHandler).Methods("POST")
	b.router.HandleFunc("/post.html", b.viewPostHandler).Methods("GET")
	b.router.HandleFunc("/edit.html", b.editPostHandler).Methods("GET")
	b.router.HandleFunc("/save", b.savePostHandler).Methods("POST")
	b.router.HandleFunc("/delete", b.deletePostHandler).Methods("GET")
	b.router.HandleFunc("/{custom}", b.customHandler).Methods("GET")

	b.router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(wd+"/skins/"+b.theme+"/assets"))))
	b.router.PathPrefix("/common_assets/").Handler(http.StripPrefix("/common_assets/", http.FileServer(http.Dir(wd+"/common_assets"))))

	http.ListenAndServe(":"+b.port, b.router)
}

var editor = `
<div id="anchor">
</div>
<div id="post-editor" title="Post editor" >
<div style="height:100%">
  <form id="form" method="POST" action="save" style="height:100%">
    <input type="text" name="title" id="input_title" value="{{.Title}}" placeholder="Title" style="width:100%; margin:5px 0px"/>

    <input type="text" name="short_description" id="input_short_description" value="{{.ShortDescription}}" placeholder="Short Decription" style="width:100%; margin-bottom:5px"/>
    <div style="margin-bottom:5px">
      <button id="btnMD_bold" type="button" class="btn btn-default" aria-label="Left Align" title="Bold">
        <span class="fa fa-bold fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_italic" type="button" class="btn btn-default" aria-label="Left Align" title="Italic">
        <span class="fa fa-italic fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_H1" type="button" class="btn btn-default" aria-label="Left Align" title="Header 1">
        <span class="fa fa-header fa-fw" aria-hidden="true"/>1
      </button>
      <button id="btnMD_H2" type="button" class="btn btn-default" aria-label="Left Align" title="Header 2">
        <span class="fa fa-header fa-fw" aria-hidden="true"/>2
      </button>
      <button id="btnMD_H3" type="button" class="btn btn-default" aria-label="Left Align" title="Header 3">
        <span class="fa fa-header fa-fw" aria-hidden="true"/>3
      </button>
      <button id="btnMD_image" type="button" class="btn btn-default" aria-label="Left Align" title="Image">
        <span class="fa fa-photo fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_hyperlink" type="button" class="btn btn-default" aria-label="Left Align" title="Hyperlink">
        <span class="fa fa-link fa-fw" aria-hidden="true"/>
      </button>


      <button id="btnMD_quote" type="button" class="btn btn-default" aria-label="Left Align" title="Quote">
        <span class="fa fa-quote-left fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_list" type="button" class="btn btn-default" aria-label="Left Align" title="List">
        <span class="fa fa-list-ul fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_orderedlist" type="button" class="btn btn-default" aria-label="Left Align" title="Ordered List">
        <span class="fa fa-list-ol fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_code" type="button" class="btn btn-default" aria-label="Left Align" title="Code">
        <span class="fa fa-code fa-fw" aria-hidden="true"/>
      </button>
      <button id="btnMD_linebreak" type="button" class="btn btn-default" aria-label="Left Align" title="Horizontal Line">
        <span class="fa fa-ellipsis-h fa-fw" aria-hidden="true"/>
      </button>

    </div>
    <textarea name="src_content" id="src" style="width:100%; height:90%; font-family:Consolas,Monaco,Lucida Console,Liberation Mono,DejaVu Sans Mono,Bitstream Vera Sans Mono,Courier New, monospace;">{{.Content}}</textarea>
    <input type="hidden" name="id" value="{{.ID}}"/>
    <input type="hidden" id="html_content" name="html_content" value=""/>
  </form>
</div>
</div>`
