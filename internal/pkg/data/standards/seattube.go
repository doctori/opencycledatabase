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

const stCollection = "seattubes"

// SeatTube will define the SeatTube (not the SeatTubeset) standard
type SeatTube struct {
	Standard `formType:"-"`
	//Diameter in mm
	Diameter float32 `formType:"int" formUnit:"mm"`
}

func NewSeatTube() *SeatTube {
	st := new(SeatTube)
	st.Type = "SeatTube"
	st.CompatibleTypes = []string{
		"Frame",
	}
	handledStandard[st.Type] = stCollection
	return st
}

// Init will setup a few fields that are immutable to the struct
func (st *SeatTube) Init() {
	st.Type = "SeatTube"
	st.CompatibleTypes = []string{
		"Frame",
	}
	st.ID = primitive.NewObjectID()
}

func (st *SeatTube) GetCompatibleTypes() []string {
	return st.CompatibleTypes
}

// Get SeatTube return the requests SeatTube Standards ID
func (st *SeatTube) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return st.Standard.Get(db, values, id, st, adj)
}

// Delete SeatTube delete the requested SeatTube standard ID
func (st *SeatTube) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return st.Standard.Delete(db, values, id, st)
}

// Post SeatTube delete the requested SeatTube standard ID
func (st *SeatTube) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return st.Standard.Post(db, values, request, id, adj, st)
}

// Put SeatTube delete the requested SeatTube standard ID
func (st *SeatTube) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return st.Standard.Put(db, values, body, st)
}

// Save SeatTube will register the SeatTube into the database
func (st *SeatTube) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[st.GetType()]
	col := db.Collection(collectionName)
	if st.ID == primitive.NilObjectID {
		st.ID = primitive.NewObjectID()
		log.Printf("Object of type %s is new inserting it into collection %s", st.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), st)
		log.Print(res)
		return
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": st.ID}
	_, err = col.UpdateOne(context.TODO(), filter, st, opts)
	return
}
