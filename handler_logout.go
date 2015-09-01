package main

import (
	"fmt"
	"log"
	"net/http"
)

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	session, err := store.Get(request, "blog-session")
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

	session.Values["admin-logged"] = false
	err = session.Save(request, response)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	http.Redirect(response, request, "/admin", 302)
}
