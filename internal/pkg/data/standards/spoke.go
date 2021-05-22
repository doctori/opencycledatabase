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

const sCollection = "spokes"

type Spoke struct {
	Standard `gorm:"embedded" formType:"-"`
	// Length is the length of the spoke
	Length int16 `formType:"int" formUnit:"mm"`
}

func NewSpoke() *Spoke {
	s := new(Spoke)
	s.Init()
	handledStandard[s.Type] = sCollection
	return s
}

// Init will setup a few fields that are immutable to the struct
func (s *Spoke) Init() {
	s.Type = "Spoke"
	s.CompatibleTypes = []string{
		"Hub",
		"Rim",
	}
	s.ID = primitive.NewObjectID()
}

// Get Spoke return the requests Spoke Standards ID
func (s *Spoke) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return s.Standard.Get(db, values, id, s, adj)
}

// Delete Spoke delete the requested Spoke standard ID
func (s *Spoke) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return s.Standard.Delete(db, values, id, s)
}

// Post Spoke delete the requested Spoke standard ID
func (s *Spoke) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return s.Standard.Post(db, values, request, id, adj, s)
}

// Put Spoke delete the requested Spoke standard ID
func (s *Spoke) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return s.Standard.Put(db, values, body, s)
}

// Save Spoke will register the crank into the database
func (s *Spoke) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[s.GetType()]
	col := db.Collection(collectionName)
	if s.ID == primitive.NilObjectID {
		s.ID = primitive.NewObjectID()
		log.Printf("Object of type %s is new inserting it into collection %s", s.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), s)
		log.Print(res)
		return
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": s.ID}
	_, err = col.UpdateOne(context.TODO(), filter, s, opts)
	return
}
