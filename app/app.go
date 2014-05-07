package mailpin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zachlatta/go-mailpin/controller/email"
	"github.com/zachlatta/go-mailpin/model"
	"github.com/zachlatta/go-mailpin/view"

	"appengine"
)

func init() {
	r := mux.NewRouter()

	r.Handle("/", model.AppHandler(root)).Methods("GET")
	r.Handle("/{id}", model.AppHandler(viewEmail)).Methods("GET")
	r.Handle("/_ah/mail/", model.AppHandler(email.IncomingHandler)).
		Methods("POST")

	http.Handle("/", r)
}

func root(w http.ResponseWriter, r *http.Request) *model.AppError {
	return view.RenderTemplate(w, "index", nil)
}

func viewEmail(w http.ResponseWriter, r *http.Request) *model.AppError {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	id := vars["id"]

	page, err := model.GetPage(c, id)
	if err != nil {
		return &model.AppError{err, "Page not found", http.StatusNotFound}
	}

	w.Write(page.Body)
	return nil
}
