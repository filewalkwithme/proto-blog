package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type loginPage struct {
	BlogTitle       string
	BlogDescription string
	ShowError       bool
	Error           string
}

//create function to log and handle error messages
func loginHandler(response http.ResponseWriter, request *http.Request) {
	session, err := store.Get(request, "blog-session")
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	if session.Values["admin-logged"] == true {
		http.Redirect(response, request, "/", 302)
		return
	}

	err = request.ParseForm()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	pUsername := request.FormValue("username")
	pPassword := request.FormValue("password")

	envPassword := os.Getenv("blog_password_" + pUsername)
	if len(envPassword) > 0 {
		if authorUsername == pUsername && envPassword == pPassword {
			session.Values["admin-logged"] = true

			err = session.Save(request, response)
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				log.Printf("%v \n", err)
				fmt.Fprintf(response, "%v \n", err)
				return
			}

			http.Redirect(response, request, "/", 302)
			return
		}
	}

	session.AddFlash("Username/password incorrect!")
	err = session.Save(request, response)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	http.Redirect(response, request, "/admin", 302)
}

// login page
func loginPageHandler(response http.ResponseWriter, request *http.Request) {
	session, err := store.Get(request, "blog-session")
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	if session.Values["admin-logged"] == true {
		http.Redirect(response, request, "/", 302)
		return
	}

	flash := ""

	flashes := session.Flashes()
	showError := len(flashes) > 0
	if showError {
		flash = flashes[0].(string)
	}
	err = session.Save(request, response)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	var page = loginPage{BlogTitle: blogTitle, BlogDescription: blogDescription, ShowError: showError, Error: flash}

	bufPage, err := ioutil.ReadFile("skins/" + theme + "/login.html")
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}
	t := template.Must(template.New("loginPage").Parse(string(bufPage)))

	t.Execute(response, page)
}
