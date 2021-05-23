package standards

import (
	"io"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const sCollection = "spokes"

type Spoke struct {
	Standard `formType:"-"`
	// Length is the length of the spoke
	Length int16 `formType:"int" formUnit:"mm"`
}

func NewSpoke() *Spoke {
	s := new(Spoke)
	s.Init()
	handledStandard[s.Type] = sCollection
	return s
}

// Init will setup a few fields that are immutable to the struct
func (s *Spoke) Init() {
	s.Type = "Spoke"
	s.CompatibleTypes = []string{
		"Hub",
		"Rim",
	}
	s.ID = primitive.NewObjectID()
}

func (s *Spoke) GetCompatibleTypes() []string {
	return s.CompatibleTypes
}

// Get Spoke return the requests Spoke Standards ID
func (s *Spoke) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return s.Standard.Get(db, values, id, s, adj)
}

// Delete Spoke delete the requested Spoke standard ID
func (s *Spoke) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return s.Standard.Delete(db, values, id, s)
}

// Post Spoke delete the requested Spoke standard ID
func (s *Spoke) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return s.Standard.Post(db, values, request, id, adj, s)
}

// Put Spoke delete the requested Spoke standard ID
func (s *Spoke) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return s.Standard.Put(db, values, body, s)
}
