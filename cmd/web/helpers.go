package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"suryanshmak.net/snippetBox/pkg/models"
	"suryanshmak.net/snippetBox/pkg/models/forms"
)

type templateData struct {
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	CurrentYear       int
	Form              *forms.Form
	AuthenticatedUser string
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
// 	if td == nil {
// 		td = &templateData{}
// 	}
// 	td.CurrentYear = time.Now().Year()
// 	td.AuthenticatedUser = app.authenticatedUser(r)
// 	return td
// }

// func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
// 	ts, ok := app.templateCache[name]
// 	if !ok {
// 		app.serverError(w, fmt.Errorf("template %s does not exist", name))
// 		return
// 	}

// 	buf := new(bytes.Buffer)

// 	err := ts.Execute(buf, app.addDefaultData(td, r))
// 	if err != nil {
// 		app.serverError(w, err)
// 		return
// 	}

// 	buf.WriteTo(w)
// }

// func (app *application) authenticatedUser(r *http.Request) string {
// 	cookie, err := r.Cookie("email")
// 	if err != nil {
// 		return ""
// 	}
// 	return cookie.Value
// }