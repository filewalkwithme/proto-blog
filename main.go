package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()

func main() {
	router.HandleFunc("/", indexPageHandler).Methods("GET")

	wd, _ := os.Getwd()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(wd+"/assets"))))

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
