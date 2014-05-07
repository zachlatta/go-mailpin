package mailpin

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"appengine"
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
	r.Handle("/", appHandler(HandleRoot))
	http.Handle("/", r)
}

func HandleRoot(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Fprintln(w, "Hello, World!")
	return nil
}
