package mailpin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/gorilla/mux"
	"github.com/zachlatta/go-mailpin/model"

	"appengine"
	"appengine/datastore"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.ERror.
		c := appengine.NewContext(r)
		c.Errorf("%v", e.Error)
		http.Error(w, e.Message, e.Code)
	}
}

func init() {
	r := mux.NewRouter()
	r.Handle("/", appHandler(root)).Methods("GET")
	r.Handle("/{id}", appHandler(viewEmail)).Methods("GET")
	r.Handle("/_ah/mail/", appHandler(incomingMail)).Methods("POST")
	http.Handle("/", r)
}

func root(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Fprintln(w, `go-mailpin is an open source clone of http://mailp.in.

Mail to p@go-mailpin.appspotmail.com. Get a short sharable URL.
  `)
	return nil
}

func viewEmail(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	id := vars["id"]

	page, err := model.GetPage(c, id)
	if err != nil {
		return &appError{err, "Page not found", http.StatusNotFound}
	}

	w.Write(page.Body)
	return nil
}

func incomingMail(w http.ResponseWriter, r *http.Request) *appError {
	c := appengine.NewContext(r)
	defer r.Body.Close()

	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		return &appError{err, "Error parsing email", http.StatusBadRequest}
	}

	bytes, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		return &appError{err, "Internal server error",
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
		return &appError{err, "Error storing page", http.StatusInternalServerError}
	}
	return nil
}
