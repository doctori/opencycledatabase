package data

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const brandCollection = "brands"

// Brand hold the brand defintion
type Brand struct {
	ID                 primitive.ObjectID `bson:"_id"`
	Name               string             `bson:"name"`
	Description        string             `bson:"description"`
	Image              primitive.ObjectID `bson:"image"`
	CreationYear       int                `bson:"creationYear"`
	EndYear            int                `bson:"endYear"`
	Country            string             `bson:"country"`
	WikiHref           string             `bson:"wikiHref"`
	Href               string             `bson:"href"`
	PutNotSupported    `bson:"-"`
	DeleteNotSupported `bson:"-"`
}

// Get will return the asked Brand
func (Brand) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	ctx := context.TODO()
	col := db.Collection(brandCollection)
	log.Printf("ID is %s", id)
	if id == primitive.NilObjectID {
		var brands []Brand

		cursor, err := col.Find(ctx, bson.M{})
		if err != nil {
			return 500, err.Error()
		}
		defer cursor.Close(ctx)
		err = cursor.All(context.TODO(), &brands)
		if err != nil {
			log.Print(err.Error())
			return 500, err.Error()
		}
		log.Print(brands)
		return 200, brands
	}
	// we have an ID return one result
	var brand Brand
	if err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&brand); err != nil {
		log.Printf("Could not decode brand fetch because %s", err.Error())
		return 404, err.Error()
	}

	return 200, brand
}

// Post will save the brand
func (Brand) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	body := request.Body
	var brand Brand
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&brand)
	if err != nil {
		log.Println(err)
		return 500, err.Error()
	}
	log.Println(brand)
	brand = brand.save(db)
	return 200, brand
}

func (b Brand) save(db *mongo.Database) Brand {
	col := db.Collection(brandCollection)
	if b.ID == primitive.NilObjectID {
		b.ID = primitive.NewObjectID()
		col.InsertOne(context.TODO(), &b)
		return b
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": b.ID}
	col.UpdateOne(context.TODO(), filter, b, opts)
	return b
}

// Delete the bike on DB given an bike ID -int-
func (Brand) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	if id != primitive.NilObjectID {
		log.Print("Will Delete the ID : ")
		log.Println(id)
		col := db.Collection(brandCollection)
		deleteResult, err := col.DeleteOne(context.TODO(), bson.M{"_id": id})
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
