package standards

import (
	"io"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const crCollection = "chainrings"

// ChainRing will define the Chainrings standard
type ChainRing struct {
	Standard `formType:"-"`
	// BoltCircleDiameter of the chainring (ref : https://www.sheldonbrown.com/gloss_bo-z.html#bcd)
	BoltCircleDiameter float32 `formType:"int" formUnit:"mm"`
	// BoltsNumber hold the number of bolt on the chainring
	BoltsNumber int `formType:"int" formUnit:"count"`
	// IsIntegrated is true if the chainring is soldered to the crank
	IsIntegrated bool `formType:"bool"`
	// IsDirectlyMounted weither the chainring is a direct mount chainring or not
	IsDirectlyMounted bool `formType:"bool"`
	// Teeth is the number of teeth a chainring has (0-255)
	Teeth uint8 `formType:"int" formUnit:"count"`
}

// NewChainRing return a ChainRing empty object with some predefined fields
func NewChainRing() *ChainRing {
	cr := new(ChainRing)
	cr.Init()
	handledStandard[cr.Type] = crCollection
	return cr
}

// Init will setup a few fields that are immutable to the struct
func (cr *ChainRing) Init() {
	cr.Type = "ChainRing"
	cr.CompatibleTypes = []string{
		"Chain",
		"Crank",
	}
	cr.ID = primitive.NewObjectID()
}
func (cr *ChainRing) GetCompatibleTypes() []string {
	return cr.CompatibleTypes
}

// Get ChainRing return the requests ChainRing Standards ID
func (cr *ChainRing) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return cr.Standard.Get(db, values, id, cr, adj)
}

// Delete ChainRing delete the requested ChainRing standard ID
func (cr *ChainRing) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return cr.Standard.Delete(db, values, id, cr)
}

// Post ChainRing delete the requested ChainRing standard ID
func (cr *ChainRing) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return cr.Standard.Post(db, values, request, id, adj, cr)
}

// Put ChainRing delete the requested ChainRing standard ID
func (cr *ChainRing) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return cr.Standard.Put(db, values, body, cr)
}
