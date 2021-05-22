package standards

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const tCollection = "threads"

// Thread defines the the standard of a thread used
// todo : https://en.wikipedia.org/wiki/Screw_thread
type Thread struct {
	Standard      `gorm:"embedded" formType:"-"`
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

// Save Thread will register the t into the database
func (t *Thread) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[t.GetType()]
	col := db.Collection(collectionName)
	if t.ID == primitive.NilObjectID {
		t.ID = primitive.NewObjectID()
		log.Printf("Object of type %s is new inserting it into collection %s", t.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), t)
		log.Print(res)
		return
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": t.ID}
	_, err = col.UpdateOne(context.TODO(), filter, t, opts)
	return
}
