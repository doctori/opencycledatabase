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

const cCollection = "cranks"

// Crank will define the Crank (not the crankset) standard
type Crank struct {
	Standard `formType:"-"`
	// Length of the crank (cm)
	Length float32 `formType:"int" formUnit:"cm"`
}

func NewCrank() *Crank {
	c := new(Crank)
	c.Init()
	handledStandard[c.Type] = cCollection
	return c
}

// Init will setup a few fields that are immutable to the struct
func (c *Crank) Init() {
	c.Type = "Crank"
	c.CompatibleTypes = []string{
		"BottomBracket",
		"ChainRing",
	}
	c.ID = primitive.NewObjectID()
}
func (c *Crank) GetCompatibleTypes() []string {
	return c.CompatibleTypes
}

// Get Crank return the requests Crank Standards ID
func (c *Crank) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return c.Standard.Get(db, values, id, c, adj)
}

// Delete Crank delete the requested Crank standard ID
func (c *Crank) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return c.Standard.Delete(db, values, id, c)
}

// Post Crank delete the requested Crank standard ID
func (c *Crank) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return c.Standard.Post(db, values, request, id, adj, c)
}

// Put Crank delete the requested Crank standard ID
func (c *Crank) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return c.Standard.Put(db, values, body, c)
}

// Save Crank will register the crank into the database
func (c *Crank) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[c.GetType()]
	col := db.Collection(collectionName)
	if c.ID == primitive.NilObjectID {
		c.ID = primitive.NewObjectID()
		log.Printf("Object of type %s is new inserting it into collection %s", c.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), c)
		log.Print(res)
		return
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": c.ID}
	_, err = col.UpdateOne(context.TODO(), filter, c, opts)
	return
}
