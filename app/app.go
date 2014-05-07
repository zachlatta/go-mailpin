package mailpin

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/gorilla/mux"
	"github.com/zachlatta/go-mailpin/model"
	"github.com/zachlatta/go-mailpin/view"

	"appengine"
	"appengine/datastore"
)

type AppHandler func(http.ResponseWriter, *http.Request) *model.AppError

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *model.AppError, not os.ERror.
		c := appengine.NewContext(r)
		c.Errorf("%v", e.Error)
		http.Error(w, e.Message, e.Code)
	}
}

func init() {
	r := mux.NewRouter()
	r.Handle("/", AppHandler(root)).Methods("GET")
	r.Handle("/{id}", AppHandler(viewEmail)).Methods("GET")
	r.Handle("/_ah/mail/", AppHandler(incomingMail)).Methods("POST")
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

func incomingMail(w http.ResponseWriter, r *http.Request) *model.AppError {
	c := appengine.NewContext(r)
	defer r.Body.Close()

	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		return &model.AppError{err, "Error parsing email", http.StatusBadRequest}
	}

	bytes, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		return &model.AppError{err, "Internal server error",
			http.StatusInternalServerError}
	}

	hasher := md5.New()
	hasher.Write(bytes)
	id := hex.EncodeToString(hasher.Sum(nil))

	page := model.Page{
		Subject: msg.Header["Subject"][0],
		Body:    bytes,
	}

	key := model.NewPageKey(c, id)
	if _, err := datastore.Put(c, key, &page); err != nil {
		return &model.AppError{err, "Error storing page", http.StatusInternalServerError}
	}
	return nil
}
