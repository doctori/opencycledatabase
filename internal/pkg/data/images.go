package data

import (
	"context"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const imageCollection = "images"

// Image is the description and the pointer to the image
type Image struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          string             `bson:"name"`
	Path          string             `bson:"path"`
	Type          string             `bson:"type"`
	ContentType   string             `bson:"contentType"`
	ContentLength int64              `bson:"contentLength"`
	Content       []byte             `bson:"-"`
	PutNotSupported
	DeleteNotSupported
}

// Get will return the image (content type is image/jpeg)
// TODO : make the content detection working
func (Image) Get(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{}) {
	var img Image
	if id == primitive.NilObjectID {
		return 200, img.List(db)
	}
	col := db.Collection(imageCollection)
	result := col.FindOne(context.TODO(), bson.M{"_id": id})
	err := result.Decode(&img)
	if err != nil {
		return 404, "Image Not Found"
	}
	file, err := os.Open(img.Path)
	log.Printf("Serving : %v", img.Path)
	defer file.Close()
	if err != nil {
		log.Println("Could Not Find the Image")
		return 404, "File Not Found"
	}
	// Detect Mime Type
	buff := make([]byte, 512)
	_, err = file.ReadAt(buff, 0)
	if err != nil {
		log.Panic(err)
	}
	img.ContentType = http.DetectContentType(buff)
	log.Printf("ContentType is : %s", img.ContentType)
	return 200, img

}

// Post saves the images to the "upload" directory (shouldn't it go to S3 ? )
func (Image) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{}) {
	uploadFolder := "upload/"
	mediaType, params, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	checkErr(err, "Could Not Determine the Content Type")
	img := Image{}
	if strings.HasPrefix(mediaType, "multipart/") {
		log.Println("YOU ARE POSTING AN IMAGE")
		multipartReader := multipart.NewReader(request.Body, params["boundary"])
		// Should Build buffer and Write it at the end (loop thour the nextpart !)
		partReader, err := multipartReader.NextPart()
		if err != nil {
			log.Printf("Whololo could not read the multiplart form submission because %s", err)
		}
		fileName := partReader.FileName()
		filePath := uploadFolder + fileName
		file, err := os.Create(filePath)
		if err != nil {
			log.Println(err)
			return 500, err
		}
		checkErr(err, "Could Not Get the Next Part")
		if _, err := io.Copy(file, partReader); err != nil {
			log.Fatal(err)
			return 500, err
		}
		// Let's record this image and return it to our client
		img = Image{Name: fileName, Path: filePath}
		img.save(db)
		log.Printf("Posted : %#v", img)
		return 200, img
	}

	return 200, img
}
func (Image) List(db *mongo.Database) (images []Image) {
	col := db.Collection(imageCollection)
	cursor, err := col.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Panicf("Could not make a list of images %s", err.Error())
	}
	cursor.All(context.TODO(), &images)
	return
}
func (i *Image) save(db *mongo.Database) (err error) {
	col := db.Collection(imageCollection)
	if i.ID == primitive.NilObjectID {
		i.ID = primitive.NewObjectID()
		_, err = col.InsertOne(context.TODO(), i)
		return
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": i.ID}
	_, err = col.UpdateOne(context.TODO(), filter, i, opts)

	return
}
