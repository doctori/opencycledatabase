package data

import (
	"io"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

type (
	// GetNotSupported default response when you cannot Get the resource
	GetNotSupported struct{}
	// PostNotSupported default response when you cannot Post the resource
	PostNotSupported struct{}
	//PutNotSupported default response when you cannot Get the resource
	PutNotSupported struct{}
	// DeleteNotSupported default response when you cannot Get the resource
	DeleteNotSupported struct{}
)

// Get returns a 405
func (GetNotSupported) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return 405, ""
}

// Post returns a 405
func (PostNotSupported) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return 405, ""
}

// Put returns a 405
func (PutNotSupported) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return 405, ""
}

// Delete returns a 405
func (DeleteNotSupported) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return 405, ""
}
