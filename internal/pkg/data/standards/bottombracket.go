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

const bbCollection = "bottombrackets"

// BottomBracket will define the bottom bracket standard
type BottomBracket struct {
	Standard `gorm:"embedded" formType:"-"`
	// Thread definition (if needed)
	ThreadID int    `json:"-" fromType:"-"`
	Thread   Thread `formType:"nested"`
	// IsThreaded : true if  it's a threaded bottom bracket
	IsThreaded bool `json:"isThreaded" formType:"bool"`
	// IsPressFit : true if it's a pressfit bottom bracket
	// (can't be true with isThreaded)
	IsPressFit bool `json:"isPressFit" formType:"bool"`
	// the inside width of the shell in mm
	ShellWidth float32 `json:"shellWidth" formType:"int" formUnit:"mm"`
	// External diameter in mm
	ExternalWidth float32 `json:"externalWidth" formType:"int" formUnit:"mm"`
}

func NewBottomBracket() *BottomBracket {
	bb := new(BottomBracket)
	bb.Init()
	// TODO : the other ones
	if len(handledStandard) == 0 {
		handledStandard = make(map[string]string)
	}
	handledStandard[bb.Type] = bbCollection
	return bb
}

// Init will setup a few fields that are immutable to the struct
func (bb *BottomBracket) Init() {
	bb.Type = "BottomBracket"
	bb.CompatibleTypes = []string{
		"Crank",
		"Frame",
	}
	bb.ID = primitive.NewObjectID()
}

// Get BottomBracket return the requests BottomBracket Standards ID
func (bb *BottomBracket) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return bb.Standard.Get(db, values, id, bb, adj)
}

// Delete BottomBracket delete the requested BottomBracket standard ID
func (bb *BottomBracket) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return bb.Standard.Delete(db, values, id, bb)
}

// Post BottomBracket delete the requested BottomBracket standard ID
func (bb *BottomBracket) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return bb.Standard.Post(db, values, request, id, adj, bb)
}

// Put BottomBracket delete the requested BottomBracket standard ID
func (bb *BottomBracket) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return bb.Standard.Put(db, values, body, bb)
}

// Save BottomBracket will register the BB into the database
func (bb *BottomBracket) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[bb.GetType()]
	col := db.Collection(collectionName)
	if bb.ID == primitive.NilObjectID {
		bb.Init()
		log.Printf("Object of type %s is new inserting it into collection %s", bb.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), bb)
		log.Print(res)
		return
	}
	upsert := true
	opts := options.FindOneAndUpdateOptions{
		Upsert: &upsert,
	}
	filter := bson.M{"_id": bb.ID}
	col.FindOneAndUpdate(context.TODO(), filter, bb, &opts)
	return
}
