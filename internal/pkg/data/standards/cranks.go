package standards

import (
	"io"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const cCollection = "cranks"

// Crank will define the Crank (not the crankset) standard
type Crank struct {
	Standard `formType:"-"`
	// Length of the crank (cm)
	Length float32 `formType:"int" formUnit:"cm"`
}

func NewCrank() *Crank {
	c := new(Crank)
	c.Init()
	handledStandard[c.Type] = cCollection
	return c
}

// Init will setup a few fields that are immutable to the struct
func (c *Crank) Init() {
	c.Type = "Crank"
	c.CompatibleTypes = []string{
		"BottomBracket",
		"ChainRing",
	}
	c.ID = primitive.NewObjectID()
}
func (c *Crank) GetCompatibleTypes() []string {
	return c.CompatibleTypes
}

// Get Crank return the requests Crank Standards ID
func (c *Crank) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return c.Standard.Get(db, values, id, c, adj)
}

// Delete Crank delete the requested Crank standard ID
func (c *Crank) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return c.Standard.Delete(db, values, id, c)
}

// Post Crank delete the requested Crank standard ID
func (c *Crank) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return c.Standard.Post(db, values, request, id, adj, c)
}

// Put Crank delete the requested Crank standard ID
func (c *Crank) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return c.Standard.Put(db, values, body, c)
}
