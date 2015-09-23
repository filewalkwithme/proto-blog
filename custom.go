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

type customPage struct {
	BlogTitle       string
	BlogDescription string
	Author          string
}

func (b *blog) customHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	custom := vars["custom"]

	if custom == "login.html" || custom == "edit.html" || custom == "index.html" || custom == "login.html" || custom == "post.html" {
		response.WriteHeader(http.StatusNotFound)
		log.Printf("%v \n", errors.New("Not found"))
		fmt.Fprintf(response, "%v \n", errors.New("Not found"))
		return
	}

	bufCustomPage, err := ioutil.ReadFile("skins/" + b.theme + "/" + custom)
	if err == nil {
		var page = customPage{
			BlogTitle:       b.blogTitle,
			BlogDescription: b.blogDescription,
			Author:          b.authorName}

		customPage := string(bufCustomPage)

		t := template.Must(template.New("customPage").Parse(customPage))

		t.Execute(response, page)
	} else {
		response.WriteHeader(http.StatusNotFound)
		log.Printf("%v \n", errors.New("Not found"))
		fmt.Fprintf(response, "%v \n", errors.New("Not found"))
		return
	}
}
