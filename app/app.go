package mailpin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zachlatta/go-mailpin/controller/email"
	"github.com/zachlatta/go-mailpin/controller/page"
	"github.com/zachlatta/go-mailpin/model"
	"github.com/zachlatta/go-mailpin/view"
)

func init() {
	r := mux.NewRouter()

	r.Handle("/", model.AppHandler(root)).Methods("GET")
	r.Handle("/{id}", model.AppHandler(page.View)).Methods("GET")
	r.Handle("/_ah/mail/{email}", model.AppHandler(email.IncomingHandler)).
		Methods("POST")

	http.Handle("/", r)
}

func root(w http.ResponseWriter, r *http.Request) *model.AppError {
	return view.RenderTemplate(w, "index", nil)
}
