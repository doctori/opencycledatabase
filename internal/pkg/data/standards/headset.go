package standards

import (
	"io"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const hsCollection = "headsets"

// Headset, holds the specs for a Front Derailleur
type Headset struct {
	Standard `formType:"-"`
	// SHISAnnotation the actual S.H.I.S Annotation like ZS44/28.6"	SHISAnnotation string
	SHISAnnotation string `formType:"string" bson:"SHISAnnotation"`
	// TODO : be able to list the possible values
	// FitType , is headset integrated or PressFit ?
	FitType string `formType:"string" bson:"fitType"`
	// IsThreaded
	IsThreaded bool `formType:"bool" bson:"isThreaded"`
	//CrownRaceInside Diameter
	CrownRaceInsideDiameter float32 `formType:"int" bson:"crownRaceInsideDiameter"`
	// FrameCupOutsideDiameter
	FrameCupOutsideDiameter float32 `formType:"int" bson:"frameCupOutsideDiameter"`
}

func NewHeadset() *Headset {
	hs := new(Headset)
	hs.Init()
	handledStandard[hs.Type] = hsCollection
	return hs
}

// Init will setup a few fields that are immutable to the struct
func (hs *Headset) Init() {
	hs.Type = "Headset"
	hs.CompatibleTypes = []string{
		"Frame",
		"Fork",
	}
	hs.ID = primitive.NewObjectID()
}

func (hs *Headset) GetCompatibleTypes() []string {
	return hs.CompatibleTypes
}

// Get Headset return the requests Headset Standards ID
func (hs *Headset) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	return hs.Standard.Get(db, values, id, hs, adj)
}

// Delete Headset delete the requested Headset standard ID
func (hs *Headset) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	return hs.Standard.Delete(db, values, id, hs)
}

// Post Headset delete the requested Headset standard ID
func (hs *Headset) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	return hs.Standard.Post(db, values, request, id, adj, hs)
}

// Put Headset delete the requested Headset standard ID
func (rd *Headset) Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{}) {
	return rd.Standard.Put(db, values, body, rd)
}
