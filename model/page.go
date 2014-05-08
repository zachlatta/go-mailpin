package model

import (
	"appengine"
	"appengine/datastore"
)

type Page struct {
	Subject string
	Body    string
}

func GetPage(c appengine.Context, id string) (*Page, error) {
	key := NewPageKey(c, id)

	var page Page
	if err := datastore.Get(c, key, &page); err != nil {
		return nil, err
	}

	return &page, nil
}

func NewPageKey(c appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(c, "Page", id, 0, nil)
}
