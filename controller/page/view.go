package page

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zachlatta/go-mailpin/model"
	"github.com/zachlatta/go-mailpin/view"

	"appengine"
)

func View(w http.ResponseWriter, r *http.Request) *model.AppError {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	id := vars["id"]

	page, err := model.GetPage(c, id)
	if err != nil {
		return &model.AppError{err, "Page not found", http.StatusNotFound}
	}

	return view.RenderTemplate(w, "page", page)
}
