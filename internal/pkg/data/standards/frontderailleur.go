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

const fdCollection = "frontderailleurs"

// FrontDerailleur, holds the specs for a Front Derailleur
// ref : https://www.sheldonbrown.com/front-derailers.html
type FrontDerailleur struct {
	Standard `formType:"-"`
	// TODO Specs ??
	// BrazedOn : is the Front derailleur need a brazed on fixture or a collar (Clamp On)?
	BrazedOn bool `formType:"bool" bson:"brazedOn"`
	// CollarSize
	CollarSize float32 `formType:"int" formUnit:"mm" bson:"collarSize"`
	// IsTriple if false : it's a double chainwheel only
	IsTriple bool `formType:"bool" bson:"isTriple"`
	//IsTopPull
	IsTopPull bool `formType:"bool" bson:"isTopPull"`
	//IsBottomPull
	IsBottomPull bool `formType:"bool" bson:"isBottomPull"`
	//IsSideSwing
	IsSideSwing bool `formType:"bool" bson:"isSideSwing"`
}

func NewFrontDerailleur() *FrontDerailleur {
	fd := new(FrontDerailleur)
	fd.Init()
	handledStandard[fd.Type] = fdCollection

	return fd
}

// Init will setup a few fields that are immutable to the struct
func (fd *FrontDerailleur) Init() {
	fd.Type = "FrontDerailleur"
	fd.CompatibleTypes = []string{
		"Frame",
		"RearDerailleur",
	}
	fd.ID = primitive.NewObjectID()
}

// Get FrontDerailleur return the requests FrontDerailleur Standards ID
func (fd *FrontDerailleur) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return fd.Standard.Get(db, values, id, fd, adj)
}

// Delete FrontDerailleur delete the requested FrontDerailleur standard ID
func (fd *FrontDerailleur) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return fd.Standard.Delete(db, values, id, fd)
}

// Post FrontDerailleur delete the requested FrontDerailleur standard ID
func (fd *FrontDerailleur) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return fd.Standard.Post(db, values, request, id, adj, fd)
}

// Put FrontDerailleur delete the requested FrontDerailleur standard ID
func (fd *FrontDerailleur) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return fd.Standard.Put(db, values, body, fd)
}

// Save FrontDerailleur will register the crank into the database
func (fd *FrontDerailleur) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[fd.GetType()]
	col := db.Collection(collectionName)
	if fd.ID == primitive.NilObjectID {
		fd.ID = primitive.NewObjectID()
		log.Printf("Object of type %s is new inserting it into collection %s", fd.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), fd)
		log.Print(res)
		return
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": fd.ID}
	_, err = col.UpdateOne(context.TODO(), filter, fd, opts)
	return
}
