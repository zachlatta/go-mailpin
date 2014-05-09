package email

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/zachlatta/go-mailpin/model"

	"appengine"
	"appengine/datastore"
	appMail "appengine/mail"
)

func IncomingHandler(w http.ResponseWriter, r *http.Request) *model.AppError {
	c := appengine.NewContext(r)
	defer r.Body.Close()

	msg, err := mail.ReadMessage(r.Body)
	if err != nil {
		return &model.AppError{err, "Error parsing email", http.StatusBadRequest}
	}

	addresses, err := msg.Header.AddressList("From")
	if err != nil {
		return &model.AppError{err, "Error parsing 'From' field in email",
			http.StatusBadRequest}
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
		return &model.AppError{err, "Error storing page",
			http.StatusInternalServerError}
	}

	notifyMsg := &appMail.Message{
		Sender:  "go-mailpin <noreply@go-mailpin.appspot.com>",
		To:      addressesToStrings(addresses),
		Subject: "Mailpin",
		Body:    fmt.Sprintf("https://go-mailpin.appspot.com/%s", id),
	}

	if err := appMail.Send(c, notifyMsg); err != nil {
		return &model.AppError{err, "Error sending mail",
			http.StatusInternalServerError}
	}

	return nil
}

func addressesToStrings(addresses []*mail.Address) []string {
	strings := make([]string, len(addresses))
	for i, address := range addresses {
		strings[i] = address.String()
	}
	return strings
}
