package data

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/doctori/opencycledatabase/internal/pkg/data/standards"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

const componentCollection = "components"

// Component : Generic struct to regroup most common properties
type Component struct {
	// add basic ID/Created@/Updated@/Delete@ through Gorm
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Brand       Brand              `bson:"brand"`
	Description string
	Standards   []standards.Standard `bson:"standards"`
	Images      []Image              `bson:"images"`
	Year        string               `bson:"year"`
	PutNotSupported
}

// ComponentInt : the standard COmponent interface that needs to complies to in order to be a component
type ComponentInt interface {
	GetName() string
	GetBrand() Brand
	GetDescription() string
	GetStandards() []standards.StandardInt
	GetImages() []Image
}

// Post will save the component in database
func (Component) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	body := request.Body
	fmt.Printf("Received args : \n\t %+v\n", body)
	decoder := json.NewDecoder(body)
	var component Component
	err := decoder.Decode(&component)
	if err != nil {
		log.Printf("Could not decode : %+v\nbecause %v\n", body, err)
		return 500, "Could Not Decode the Data"
	}
	log.Println(component)
	component, err = component.save(db)
	if err != nil {
		return 500, "Could Not Save the Component"
	}
	return 200, component
}

// Get return a generic Component
func (Component) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {

	page, err := strconv.ParseInt(values.Get("page"), 10, 64)
	if err != nil {
		page = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	perPage, err := strconv.ParseInt(values.Get("per_page"), 10, 64)
	if err != nil {
		perPage = defaultPerPage
	}

	var c Component
	collection := db.Collection(componentCollection)
	idx, rawFilter := bsoncore.AppendDocumentStart(nil)
	components := []Component{}
	if id != primitive.NilObjectID {
		result := collection.FindOne(context.TODO(), bson.M{"_id": id})
		if result.Err() != nil {
			return 404, fmt.Sprintf("Component not found because %s", result.Err().Error())
		}
		return 200, c
	}
	// FIXME : this is starting to look like a really wrong way of dealing with url values
	if values.Get("name") == "" && values.Get("search") == "" && values.Get("standard") == "" && values.Get("type") == "" {
		return 200, c.getAll(db, page, perPage)
	} else if values.Get("search") != "" {
		return 200, c.search(db, page, perPage, values.Get("search"))
	} else if values.Get("standard") != "" {
		bsoncore.AppendStringElement(rawFilter, "standards.name", values.Get("standard"))
		log.Printf("Standard filter is %s", values.Get("standard"))

	} else if values.Get("type") != "" {
		bsoncore.AppendStringElement(rawFilter, "type", values.Get("type"))
		log.Printf("Type filter is %s", values.Get("standard"))
	} else if values.Get("name") != "" {
		bsoncore.AppendStringElement(rawFilter, "name", values.Get("name"))
		log.Printf("Name filter is %s", values.Get("name"))
	}
	bsoncore.AppendDocumentEnd(rawFilter, idx)
	filter := bson.M{}
	err = bson.Unmarshal(rawFilter, &filter)
	if err != nil {
		return 500, fmt.Sprintf("Could not decode filter %s", err.Error())
	}
	cursor, err := collection.Find(context.TODO(), rawFilter)
	if err != nil {
		return 404, fmt.Sprintf("Component not found because %s", err.Error())
	}
	if err = cursor.All(context.TODO(), components); err != nil {
		return 500, fmt.Sprintf("Error while read mongo cursor : %s", err.Error())
	}
	return 200, components
}

func (Component) Delete(db *mongo.Database, value url.Values, id primitive.ObjectID) (int, interface{}) {

	if id == primitive.NilObjectID {
		return 400, "Component ID required"
	}
	collection := db.Collection(componentCollection)
	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return 404, fmt.Sprintf("Component not found because %s", err.Error())
	}
	log.Printf("Component is ID %d has been removed", id)
	return 201, result
}

func (c Component) search(db *mongo.Database, page int64, perPage int64, filter string) (components []Component) {
	filter = fmt.Sprintf("%%%s%%", filter)
	collection := db.Collection(componentCollection)
	skip := page * perPage
	c.find(collection,
		bson.M{"name": primitive.Regex{
			Pattern: "^" + filter,
			Options: "i"},
		},
		&options.FindOptions{
			Skip:  &skip,
			Limit: &perPage,
		},
	)

	return
}
func (c Component) getAll(db *mongo.Database, page int64, perPage int64) (components []Component) {

	//TODO : return LINKS Header with the next page and previous page
	skip := page * perPage
	collection := db.Collection(componentCollection)
	options := options.FindOptions{
		Skip:  &skip,
		Limit: &perPage,
	}
	return c.find(collection, bson.M{}, &options)

}
func (Component) find(coll *mongo.Collection, filter bson.M, options *options.FindOptions) (components []Component) {
	results, err := coll.Find(
		context.TODO(),
		filter,
		options)
	if err != nil {
		log.Printf("Could not parse find Results %s", err.Error())
		return components
	}
	err = results.All(context.TODO(), &components)
	if err != nil {
		log.Printf("Could not erad cursor %s ", err.Error())
	}
	return components
}
func (c Component) save(db *mongo.Database) (Component, error) {
	opts := options.Update().SetUpsert(true)
	col := db.Collection(brandCollection)
	filter := bson.M{"_id": c.ID}
	_, err := col.UpdateOne(context.TODO(), filter, c, opts)

	return c, err
}
