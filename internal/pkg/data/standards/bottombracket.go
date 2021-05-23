package standards

import (
	"io"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bbCollection = "bottombrackets"

// BottomBracket will define the bottom bracket standard
type BottomBracket struct {
	Standard `formType:"-"`
	// Thread definition (if needed)
	Thread Thread `formType:"nested"`
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

func (bb *BottomBracket) GetCompatibleTypes() []string {
	return bb.CompatibleTypes
}
