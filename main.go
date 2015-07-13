package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var router = mux.NewRouter()

//DB is the global DB object
var DB *sql.DB

func init() {
	_, errCheckExists := os.Stat("./blog.db")

	db, err := sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	if os.IsNotExist(errCheckExists) {
		fmt.Printf("File blog.db not exists. \nCreating initial database.\n")

		sqlCreateDB := `create table post (id integer not null primary key, post text, date datetime);`

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

	wd, _ := os.Getwd()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(wd+"/assets"))))

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
