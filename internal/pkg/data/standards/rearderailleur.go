package standards

import (
	"context"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const rdCollection = "rearderailleurs"

type RearDerailleur struct {
	Standard `formType:"-"`

	// CageLength hold the length of the cage in mm (is this a standard defining key ?)
	CageLength float32 `json:"CageLength" formType:"int"  formUnit:"mm"`
	// IsDirectMount ? ref : https://wheelsmfg.com/blog/standard-mount-vs-direct-mount-derailleur-hangers.html
	IsDirectMount bool `json:"IsDirectMount" formType:"bool"`
	// IsShortCage will say if it's short or if its long (do we need a "longcage" option ? )
	IsShortCage bool `json:"IsShortCage" formType:"bool"`
	IsLongCage  bool `json:"IsLongCage" formType:"bool"`
}

func NewRearDerailleur() *RearDerailleur {
	rd := new(RearDerailleur)
	rd.Init()
	handledStandard[rd.Type] = rdCollection
	return rd
}

// Init will setup a few fields that are immutable to the struct
func (rd *RearDerailleur) Init() {
	rd.Type = "RearDerailleur"
	rd.CompatibleTypes = []string{
		"Cassette",
		"Frame",
	}
	rd.ID = primitive.NewObjectID()
}

func (rd *RearDerailleur) GetCompatibleTypes() []string {
	return rd.CompatibleTypes
}

// Get RearDerailleur return the requests RearDerailleur Standards ID
func (rd *RearDerailleur) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return rd.Standard.Get(db, values, id, rd, adj)
}

// Delete RearDerailleur delete the requested RearDerailleur standard ID
func (rd *RearDerailleur) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return rd.Standard.Delete(db, values, id, rd)
}

// Post RearDerailleur delete the requested RearDerailleur standard ID
func (rd *RearDerailleur) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return rd.Standard.Post(db, values, request, id, adj, rd)
}

// Put RearDerailleur delete the requested RearDerailleur standard ID
func (rd *RearDerailleur) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return rd.Standard.Put(db, values, body, rd)
}

// Save RearDerailleur will register the crank into the database
func (rd *RearDerailleur) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[rd.GetType()]
	col := db.Collection(collectionName)
	if rd.ID == primitive.NilObjectID {
		rd.ID = primitive.NewObjectID()
		log.Printf("Object of type %s is new inserting it into collection %s", rd.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), rd)
		log.Print(res)
		return
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": rd.ID}
	_, err = col.UpdateOne(context.TODO(), filter, rd, opts)
	return
}
