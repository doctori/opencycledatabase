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

const wheelCollection = "wheels"

type Wheel struct {
	Standard `gorm:"embedded"`
	// Diameter : the Diameter of the wheel (let's say mm)
	Diameter int16 `formType:"int" formUnit:"mm"`
	// Should we include the subStandard ? Hub/Spoke/Rim ?
}

func NewWheel() *Wheel {

	w := new(Wheel)
	w.Type = "Wheel"
	w.CompatibleTypes = []string{"Fork", "Tire"}
	handledStandard[w.Type] = wheelCollection
	return w
}

// Init will setup a few fields that are immutable to the struct
func (w *Wheel) Init() {
	w.Type = "Wheel"
	w.CompatibleTypes = []string{
		"Fork",
		"Tire",
	}
	w.ID = primitive.NewObjectID()
}

// Get Wheel return the requests Wheel Standards ID
func (w *Wheel) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return w.Standard.Get(db, values, id, w, adj)
}

// Delete Wheel delete the requested Wheel standard ID
func (w *Wheel) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return w.Standard.Delete(db, values, id, w)
}

// Post Wheel delete the requested Wheel standard ID
func (w *Wheel) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return w.Standard.Post(db, values, request, id, adj, w)
}

// Put Wheel delete the requested Wheel standard ID
func (w *Wheel) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return w.Standard.Put(db, values, body, w)
}

// Save Wheel will register the crank into the database
func (w *Wheel) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[w.GetType()]
	col := db.Collection(collectionName)
	if w.ID == primitive.NilObjectID {
		w.ID = primitive.NewObjectID()
		log.Printf("Object of type %s is new inserting it into collection %s", w.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), w)
		log.Print(res)
		return
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": w.ID}
	_, err = col.UpdateOne(context.TODO(), filter, w, opts)
	return
}
