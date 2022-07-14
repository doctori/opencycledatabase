package utils

import (
	"io"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
func (GetNotSupported) Get(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return 405, ""
}

// Post returns a 405
func (PostNotSupported) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return 405, ""
}

// Put returns a 405
func (PutNotSupported) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return 405, ""
}

// Delete returns a 405
func (DeleteNotSupported) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return 405, ""
}
