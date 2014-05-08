package email

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/zachlatta/go-mailpin/model"

	"appengine"
	"appengine/datastore"
)

func IncomingHandler(w http.ResponseWriter, r *http.Request) *model.AppError {
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
		Body:    string(bytes),
	}

	key := model.NewPageKey(c, id)
	if _, err := datastore.Put(c, key, &page); err != nil {
		return &model.AppError{err, "Error storing page", http.StatusInternalServerError}
	}
	return nil
}
