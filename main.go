package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	ini "github.com/Unknwon/goconfig"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var router = mux.NewRouter()

//DB is the global DB object
var DB *sql.DB
var cfg *ini.ConfigFile

func init() {
	cfg, err := ini.LoadConfigFile("blog.cfg")
	if err != nil {
		log.Fatal(err)
	}

	authorName, err := cfg.GetValue("author", "name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v \n", authorName)

	_, errCheckExists := os.Stat("./blog.db")

	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	if os.IsNotExist(errCheckExists) {
		fmt.Printf("File blog.db not exists. \nCreating initial database.\n")

		sqlCreateDB := `
		create table posts (id integer not null primary key, title text, src_content text, html_content text, date datetime);
		`

		_, err = db.Exec(sqlCreateDB)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlCreateDB)
			return
		}
		fmt.Printf("Initial database created.\n")
	}
}

func main() {
	router.HandleFunc("/", indexPageHandler).Methods("GET")
	router.HandleFunc("/post.html", postHandler).Methods("GET")
	router.HandleFunc("/edit.html", editHandler).Methods("GET")
	router.HandleFunc("/save", saveHandler).Methods("POST")
	router.HandleFunc("/delete", deleteHandler).Methods("GET")

	wd, _ := os.Getwd()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(wd+"/assets"))))

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
