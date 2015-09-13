package main

import (
	"fmt"
	"log"
	"net/http"
)

func (b *blog) logoutHandler(response http.ResponseWriter, request *http.Request) {
	session, err := b.store.Get(request, b.sessionName)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	if session.Values["admin-logged"] != true {
		http.Redirect(response, request, "/", 401)
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
