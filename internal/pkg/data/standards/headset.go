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

const hsCollection = "headsets"

// Headset, holds the specs for a Front Derailleur
type Headset struct {
	Standard `formType:"-"`
	// SHISAnnotation the actual S.H.I.S Annotation like ZS44/28.6"	SHISAnnotation string
	SHISAnnotation string `formType:"string"`
	// TODO : be able to list the possible values
	// FitType , is headset integrated or PressFit ?
	FitType string `formType:"String"`
}

func NewHeadset() *Headset {
	hs := new(Headset)
	hs.Init()
	handledStandard[hs.Type] = hsCollection
	return hs
}

// Init will setup a few fields that are immutable to the struct
func (hs *Headset) Init() {
	hs.Type = "Headset"
	hs.CompatibleTypes = []string{
		"Frame",
		"Fork",
	}
	hs.ID = primitive.NewObjectID()
}

func (hs *Headset) GetCompatibleTypes() []string {
	return hs.CompatibleTypes
}

// Get Headset return the requests Headset Standards ID
func (hs *Headset) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return hs.Standard.Get(db, values, id, hs, adj)
}

// Delete Headset delete the requested Headset standard ID
func (hs *Headset) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return hs.Standard.Delete(db, values, id, hs)
}

// Post Headset delete the requested Headset standard ID
func (hs *Headset) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return hs.Standard.Post(db, values, request, id, adj, hs)
}

// Put Headset delete the requested Headset standard ID
func (rd *Headset) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return rd.Standard.Put(db, values, body, rd)
}

// Save Headset will register the crank into the database
func (hs *Headset) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[hs.GetType()]
	col := db.Collection(collectionName)
	if hs.ID == primitive.NilObjectID {
		hs.ID = primitive.NewObjectID()
		log.Printf("Object of type %s is new inserting it into collection %s", hs.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), hs)
		log.Print(res)
		return
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": hs.ID}
	_, err = col.UpdateOne(context.TODO(), filter, hs, opts)
	return
}
