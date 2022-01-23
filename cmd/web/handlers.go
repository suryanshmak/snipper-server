package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"suryanshmak.net/snippetBox/pkg/models"
)

type CreateSnippet struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Expires string `json:"expires"`
}

type Response struct {
	ID int `json:"id"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUp struct {
	Login
	Name string `json:"name"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	json.NewEncoder(w).Encode(s)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	form := CreateSnippet{}
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires+"days")
	if err != nil {
		app.serverError(w, err)
		return
	}
	response := Response{id}
	json.NewEncoder(w).Encode(response)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	form := SignUp{}
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err == models.ErrDuplicateEmail {
		fmt.Fprint(w, err)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	http.SetCookie(w, setCookie("isSigned", "true"))
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	form := Login{}
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.users.Authenticate(form.Email, form.Password)
	if err == models.ErrInvalidCredentials {
		fmt.Fprint(w, err)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	http.SetCookie(w, setCookie("isSigned", "true"))
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, setCookie("isSigned", ""))
}

func setCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   "localhost",
		Path:     "/",
		Expires:  time.Now().UTC().Add(1000 * time.Hour),
		MaxAge:   0,
		Secure:   true,
		HttpOnly: false,
	}
}
