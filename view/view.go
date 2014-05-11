package view

import (
	"html/template"
	"net/http"

	"github.com/zachlatta/go-mailpin/model"
)

// Base directory relative to app.go where templates are stored
const tD = "view/"

var templates = template.Must(template.ParseFiles(
	tD+"index.html",
	tD+"page.html",
))

func RenderTemplate(w http.ResponseWriter, tmpl string,
	c interface{}) *model.AppError {
	err := templates.ExecuteTemplate(w, tmpl+".html", c)
	if err != nil {
		return &model.AppError{err, "Can't display webpage",
			http.StatusInternalServerError}
	}
	return nil
}
