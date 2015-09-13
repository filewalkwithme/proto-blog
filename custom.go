package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func (b *blog) customHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	custom := vars["custom"]

	bufIndexPage, err := ioutil.ReadFile("skins/" + b.theme + "/" + custom)
	if err == nil {
		indexPage := string(bufIndexPage)

		t := template.Must(template.New("indexPage").Parse(indexPage))

		t.Execute(response, nil)
	} else {
		response.WriteHeader(http.StatusNotFound)
		log.Printf("%v \n", errors.New("Not found"))
		fmt.Fprintf(response, "%v \n", errors.New("Not found"))
		return
	}
}
