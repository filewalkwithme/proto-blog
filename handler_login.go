package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type loginPage struct {
	BlogTitle       string
	BlogDescription string
	ShowError       bool
	Error           string
}

var failedAttempts = 0
var lastValidAttempt time.Time

func (b *blog) loginHandler(response http.ResponseWriter, request *http.Request) {
	if failedAttempts >= 5 && time.Since(lastValidAttempt) >= (30*time.Minute) {
		failedAttempts = 0
	}

	session, err := b.store.Get(request, b.sessionName)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}

	if failedAttempts < 5 {
		failedAttempts = failedAttempts + 1
		lastValidAttempt = time.Now()

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
			if b.authorUsername == pUsername && envPassword == pPassword {
				session.Values["admin-logged"] = true
				failedAttempts = 0

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
	} else {
		session.AddFlash("Password retry limit exceeded! Wait 30 minutes to try again.")
	}

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
func (b *blog) loginPageHandler(response http.ResponseWriter, request *http.Request) {
	session, err := b.store.Get(request, b.sessionName)
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

	var page = loginPage{
		BlogTitle:       b.blogTitle,
		BlogDescription: b.blogDescription,
		ShowError:       showError,
		Error:           flash}

	bufPage, err := ioutil.ReadFile("skins/" + b.theme + "/login.html")
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Printf("%v \n", err)
		fmt.Fprintf(response, "%v \n", err)
		return
	}
	t := template.Must(template.New("loginPage").Parse(string(bufPage)))

	t.Execute(response, page)
}
