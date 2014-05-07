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
	r.Handle("/", appHandler(root))
	http.Handle("/", r)
	http.Handle("/_ah/mail/", appHandler(incomingMail))
}

func root(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Fprintln(w, `go-mailpin is an open source clone of http://mailp.in.

Mail to p@go-mailpin.appspotmail.com. Get a short sharable URL.
  `)
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
		Body: bytes,
	}

	key := model.NewPageKey(c, id)
	if _, err := datastore.Put(c, key, &page); err != nil {
		return &appError{err, "Error storing page", http.StatusInternalServerError}
	}
	return nil
}
