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

var wd string

type blog struct {
	DB              *sql.DB
	cfg             *ini.ConfigFile
	blogTitle       string
	blogDescription string
	authorName      string
	authorUsername  string
	secret          string
	sessionName     string
	theme           string
	port            string

	store  *sessions.CookieStore
	router *mux.Router
}

func (b *blog) load(configFile string) {
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

	b.sessionName, err = cfg.GetValue("blog", "session-name")
	if err != nil {
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

	_, errCheckExists := os.Stat("./blog.db")

	db, err := sql.Open("sqlite3", "./blog.db")
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

	b.router = mux.NewRouter()
	b.store = sessions.NewCookieStore([]byte(b.secret))

	b.router.HandleFunc("/", b.indexPageHandler).Methods("GET")
	b.router.HandleFunc("/admin", b.loginPageHandler).Methods("GET")
	b.router.HandleFunc("/login", b.loginHandler).Methods("POST")
	b.router.HandleFunc("/logout", b.logoutHandler).Methods("GET")
	b.router.HandleFunc("/post.html", b.viewPostHandler).Methods("GET")
	b.router.HandleFunc("/edit.html", b.editPostHandler).Methods("GET")
	b.router.HandleFunc("/save", b.savePostHandler).Methods("POST")
	b.router.HandleFunc("/delete", b.deletePostHandler).Methods("GET")

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(wd+"/skins/"+b.theme+"/assets"))))

	http.Handle("/common_assets/", http.StripPrefix("/common_assets/", http.FileServer(http.Dir(wd+"/common_assets"))))

	http.Handle("/", b.router)

	http.ListenAndServe(":"+b.port, nil)
}

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	wd = workingDir

	b := &blog{}
	b.load("blog.cfg")
}
