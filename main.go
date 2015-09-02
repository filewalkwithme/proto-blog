package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	ini "github.com/Unknwon/goconfig"
	mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)

var router = mux.NewRouter()

//DB is the global DB object
var DB *sql.DB
var cfg *ini.ConfigFile
var blogTitle string
var blogDescription string
var authorName string
var authorUsername string
var secret string
var sessionName string
var theme string

var store = sessions.NewCookieStore([]byte(secret))

func init() {
	cfg, err := ini.LoadConfigFile("blog.cfg")
	if err != nil {
		log.Fatal(err)
	}

	authorName, err = cfg.GetValue("author", "name")
	if err != nil {
		log.Fatal(err)
	}

	authorUsername, err = cfg.GetValue("author", "username")
	if err != nil {
		log.Fatal(err)
	}

	blogTitle, err = cfg.GetValue("blog", "title")
	if err != nil {
		log.Fatal(err)
	}

	blogDescription, err = cfg.GetValue("blog", "description")
	if err != nil {
		log.Fatal(err)
	}

	secret, err = cfg.GetValue("blog", "secret")
	if err != nil {
		log.Fatal(err)
	}

	sessionName, err = cfg.GetValue("blog", "session-name")
	if err != nil {
		log.Fatal(err)
	}

	theme, err = cfg.GetValue("blog", "theme")
	if err != nil {
		log.Fatal(err)
	}

	_, errCheckExists := os.Stat("./blog.db")

	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}
	DB = db

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

func main() {
	router.HandleFunc("/", indexPageHandler).Methods("GET")
	router.HandleFunc("/admin", loginPageHandler).Methods("GET")
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("GET")
	router.HandleFunc("/post.html", viewPostHandler).Methods("GET")
	router.HandleFunc("/edit.html", editPostHandler).Methods("GET")
	router.HandleFunc("/save", savePostHandler).Methods("POST")
	router.HandleFunc("/delete", deletePostHandler).Methods("GET")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(wd+"/skins/"+theme+"/assets"))))

	http.Handle("/common_assets/", http.StripPrefix("/common_assets/", http.FileServer(http.Dir(wd+"/common_assets"))))

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
