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

const hCollection = "hubs"

// Hub holds the information on the hub
type Hub struct {
	Standard `gorm:"embedded" formType:"-"`
	// Holes : the number of spokes attached to the hub
	Holes uint8 `formType:"int"`
	// AxleType the type of axle required
	AxleType string `formType:"string"`
	// AxelLength the Length of the Axel (external)
	AxelLength float32 `formType:"int" formUnit:"mm"`
	// AxelDiameter the Diameter of the axle (external)
	AxelDiameter float32 `formType:"int" formUnit:"mm"`
	// HasDiscBrakes : does the hub support disc brakes
	HasDiscBrakes bool `formType:"bool"`
}

func NewHub() *Hub {
	h := new(Hub)
	h.Init()
	handledStandard[h.Type] = hCollection
	return h
}

// Init will setup a few fields that are immutable to the struct
func (h *Hub) Init() {
	h.Type = "Hub"
	h.CompatibleTypes = []string{
		"Spoke",
		"Disk",
		"Axle",
	}
	h.ID = primitive.NewObjectID()
}

// Get Hub return the requests Hub Standards ID
func (h *Hub) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return h.Standard.Get(db, values, id, h, adj)
}

// Delete Hub delete the requested Hub standard ID
func (h *Hub) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return h.Standard.Delete(db, values, id, h)
}

// Post Hub delete the requested Hub standard ID
func (h *Hub) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return h.Standard.Post(db, values, request, id, adj, h)
}

// Put Hub delete the requested Hub standard ID
func (h *Hub) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return h.Standard.Put(db, values, body, h)
}

// Save Hub will register the crank into the database
func (h *Hub) Save(db *mongo.Database) (err error) {
	collectionName := handledStandard[h.GetType()]
	col := db.Collection(collectionName)
	upsert := false
	if h.ID == primitive.NilObjectID {
		h.Init()
		upsert = true
	}

	opts := options.ReplaceOptions{
		Upsert: &upsert,
	}
	filter := bson.M{"_id": h.ID}
	log.Print(filter)
	res := &mongo.UpdateResult{}
	res, err = col.ReplaceOne(context.TODO(), &filter, h, &opts)
	log.Printf("%#v", res)
	if err != nil {
		log.Printf("Could not save the Hub : %s", err.Error())
		return
	}
	return
}
