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

const fCollection = "frames"

// Frame will define the Frame (not the Frameset) standard
type Frame struct {
	Standard `formType:"-"`
	// Geometry of the Frame (using https://geometrygeeks.bike/ as a reference)
	//Reach in mm
	Reach float32 `formType:"int" formUnit:"mm"`
	//Stack in mm
	Stack float32 `formType:"int" formUnit:"mm"`
	//TopTube in mm
	TopTube float32 `formType:"int" formUnit:"mm"`
	// SeatTube in mm
	SeatTube float32 `formType:"int" formUnit:"mm"`
	// HeadAngle in degres
	HeadAngle float32 `formType:"int" formUnit:"degrees"`
	// SeatAngle in degres
	SeatAngle float32 `formType:"int" formUnit:"degrees"`
	// HeadTube in mm
	HeadTube float32 `formType:"int" formUnit:"mm"`
	// Chainstay in mm
	Chainstay float32 `formType:"int" formUnit:"mm"`
	// Wheelbase in mm
	Wheelbase float32 `formType:"int" formUnit:"mm"`
	// Front Centre in mm
	FrontCentre float32 `formType:"int" formUnit:"mm"`
	// Standover in mm
	Standover float32 `formType:"int" formUnit:"mm"`
	// BBDrop in mm
	BBDrop float32 `formType:"int" formUnit:"mm"`
	// BBHeight in mm
	BBHeight float32 `formType:"int" formUnit:"mm"`
}

func NewFrame() *Frame {
	f := new(Frame)
	f.Init()
	handledStandard[f.Type] = fCollection
	return f
}

// Init will setup a few fields that are immutable to the struct
func (f *Frame) Init() {
	f.Type = "Frame"
	f.CompatibleTypes = []string{
		"FrontDerailleur",
		"RearDerailleur",
		"Fork",
		"SeatTube",
	}
	f.ID = primitive.NewObjectID()
}
func (f *Frame) GetCompatibleTypes() []string {
	return f.CompatibleTypes
}

// Get Frame return the requests Frame Standards ID
func (f *Frame) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return f.Standard.Get(db, values, id, f, adj)
}

// Delete Frame delete the requested Frame standard ID
func (f *Frame) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return f.Standard.Delete(db, values, id, f)
}

// Post Frame delete the requested Frame standard ID
func (f *Frame) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return f.Standard.Post(db, values, request, id, adj, f)
}

// Put Frame delete the requested Frame standard ID
func (f *Frame) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return f.Standard.Put(db, values, body, f)
}

// Save Frame will register the Frame into the database
func (f *Frame) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[f.GetType()]
	col := db.Collection(collectionName)
	if f.ID == primitive.NilObjectID {
		f.ID = primitive.NewObjectID()
		log.Printf("Object of type %s is new inserting it into collection %s", f.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), f)
		log.Print(res)
		return
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": f.ID}
	_, err = col.UpdateOne(context.TODO(), filter, f, opts)
	return
}
