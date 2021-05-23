package standards

import (
	"io"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hCollection = "hubs"

// Hub holds the information on the hub
type Hub struct {
	Standard `formType:"-"`
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

func (h *Hub) GetCompatibleTypes() []string {
	return h.CompatibleTypes
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
