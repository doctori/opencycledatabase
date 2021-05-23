package standards

import (
	"io"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const tCollection = "threads"

// Thread defines the the standard of a thread used
// todo : https://en.wikipedia.org/wiki/Screw_thread
type Thread struct {
	Standard      `formType:"-"`
	ThreadPerInch int16   `formType:"int" formUnit:"tpi"`
	Diameter      float32 `formType:"int" formUnit:"mm"`
	Orientation   string  `formType:"string" formUnit:"orientation"`
}

func NewThread() *Thread {
	t := new(Thread)
	t.Type = "Thread"
	t.CompatibleTypes = []string{
		"Frame",
	}
	handledStandard[t.Type] = tCollection
	return t
}

// Init will setup a few fields that are immutable to the struct
func (t *Thread) Init() {
	t.Type = "Thread"
	t.CompatibleTypes = []string{
		"Frame",
	}
	t.ID = primitive.NewObjectID()
}

func (t *Thread) GetCompatibleTypes() []string {
	return t.CompatibleTypes
}

// Get Thread return the requests Thread Standards ID
func (t *Thread) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return t.Standard.Get(db, values, id, t, adj)
}

// Delete Thread delete the requested Thread standard ID
func (t *Thread) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return t.Standard.Delete(db, values, id, t)
}

// Post Thread delete the requested Thread standard ID
func (t *Thread) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return t.Standard.Post(db, values, request, id, adj, t)
}

// Put Thread delete the requested Thread standard ID
func (t *Thread) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return t.Standard.Put(db, values, body, t)
}
