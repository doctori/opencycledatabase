package standards

import (
	"io"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const wheelCollection = "wheels"

type Wheel struct {
	Standard `formType:"-"`
	// Diameter : the Diameter of the wheel (let's say mm)
	Diameter int16 `formType:"int" formUnit:"mm" bson:"diameter"`
	// Should we include the subStandard ? Hub/Spoke/Rim ?
}

func NewWheel() *Wheel {

	w := new(Wheel)
	w.Type = "Wheel"
	w.CompatibleTypes = []string{"Fork", "Tire"}
	handledStandard[w.Type] = wheelCollection
	return w
}

// Init will setup a few fields that are immutable to the struct
func (w *Wheel) Init() {
	w.Type = "Wheel"
	w.CompatibleTypes = []string{
		"Fork",
		"Tire",
	}
	w.ID = primitive.NewObjectID()
}

func (w *Wheel) GetCompatibleTypes() []string {
	return w.CompatibleTypes
}

// Get Wheel return the requests Wheel Standards ID
func (w *Wheel) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return w.Standard.Get(db, values, id, w, adj)
}

// Delete Wheel delete the requested Wheel standard ID
func (w *Wheel) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return w.Standard.Delete(db, values, id, w)
}

// Post Wheel delete the requested Wheel standard ID
func (w *Wheel) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return w.Standard.Post(db, values, request, id, adj, w)
}

// Put Wheel delete the requested Wheel standard ID
func (w *Wheel) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return w.Standard.Put(db, values, body, w)
}
