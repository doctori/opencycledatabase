package data

import (
	"context"
	"encoding/json"
	"io"

	// import png support
	_ "image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/doctori/opencycledatabase/internal/pkg/data/standards"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultPerPage int64 = 30
const bikeCollection = "bikes"

type Resource interface {
	Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{})
	Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{})
	Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{})
	Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{})
}

// Bike contains the defintion of the bike and all its attached components
type Bike struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name              string             `bson:"name"`
	Brand             Brand              `bson:"brand,omitempty"`
	Year              string             `bson:"year,omitempty"`
	Description       string             `bson:"description,omitempty"`
	Images            []Image            `bson:"images"`
	Components        []Component
	SupportedStandard []standards.Standard
	PutNotSupported
}

// Should Return Errors !!
func (b Bike) save(db *mongo.Database) (err error) {
	col := db.Collection(bikeCollection)
	ctx := context.TODO()
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	filter := bson.M{"_id": b.ID}
	result := col.FindOneAndUpdate(ctx, filter, b, &opt)
	if result.Err() != nil {
		err = result.Err()
	}
	result.Decode(&b)
	return
}

// GetAllBikes will return all the bikes with the pagination asked
func GetAllBikes(db *mongo.Database, page string, perPage string) interface{} {
	ipage, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		ipage = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	iperPage, err := strconv.ParseInt(perPage, 10, 64)
	if err != nil {
		iperPage = defaultPerPage
	}
	var bikes []Bike
	col := db.Collection(bikeCollection)
	ctx := context.TODO()
	opt := options.FindOptions{
		Skip:  &ipage,
		Limit: &iperPage,
	}
	cursor, err := col.Find(ctx, bson.M{}, &opt)
	if err != nil {
		return bikes
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var bike Bike
		if err = cursor.Decode(&bike); err != nil {
			log.Print(err.Error())
		}
		bikes = append(bikes, bike)
	}
	return bikes
}

// Get Bike will return the request bike struct
func (Bike) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	page := values.Get("page")
	perPage := values.Get("per_page")
	/*if values.Get("name") == "" {
		return 200, GetAllBikes()
	}*/
	// Let Display All that We Have
	// Someday Pagination will be there
	if id == primitive.NilObjectID {
		return 200, GetAllBikes(db, page, perPage)
	}
	var bike Bike
	collection := db.Collection(bikeCollection)
	filter := bson.M{"_id": id}

	err := collection.FindOne(context.TODO(), filter).Decode(&bike)
	if err == mongo.ErrNoDocuments {
		return 404, "Bike not found"
	}
	return 200, bike
}

// Delete the bike on DB given an bike ID -int-
func (Bike) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	if id != primitive.NilObjectID {
		log.Print("Will Delete the ID : ")
		log.Println(id)
		col := db.Collection(bikeCollection)
		deleteResult, err := col.DeleteOne(context.TODO(), bson.M{"ID": id})
		if err == mongo.ErrNoDocuments {
			return 404, "Id Not Found"
		}
		if err != nil {
			return 500, "Bleuuharg"
		}
		return 200, deleteResult
	}
	log.Println(id)
	return 404, "NOT FOUND"
}

// Post Handle Post request from http
func (Bike) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	body := request.Body
	var bike Bike

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&bike)
	if err != nil {
		return 500, "Internal Error"
	}
	// Clean the unName Components
	components := bike.Components[:0]
	for _, component := range bike.Components {
		if component.Name != "" {
			components = append(components, component)
		}
	}
	bike.Components = components
	bike.save(db)
	return 200, bike
}

// TODO : use it
func (Bike) addComponent(db *mongo.Database, bike Bike, component Component) {
	bike.Components = append(bike.Components, component)
	bike.save(db)
}

func (Bike) getCompatibleComponents(bike Bike) []Component {
	var compatibleComponents []Component
	// Get The Components, Get their Standards

	// Get All the Components supporting those standards
	// Return Thems
	return compatibleComponents
}
