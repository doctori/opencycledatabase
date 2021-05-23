package standards

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO remove this shit from here
// this should be controlled by the API not the datamodel
const defaultPerPage int64 = 30

var handledStandard map[string]string

// StandardInt interface define all the method that a standard need to have to be a
// real standard struct
type StandardInt interface {
	GetName() string
	GetCode() string
	GetID() primitive.ObjectID
	SetID(id primitive.ObjectID)
	GetCompatibleTypes() []string
	GetType() string
	IsNul() bool
	Init()
	//GetCompatible() []StandardInt
	Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{})
	Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{})
	Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{})
	Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{})
}

// Standard define the generic common Standard properties
type Standard struct {
	ID              primitive.ObjectID   `formType:"-" bson:"_id"`
	Name            string               `formType:"string" bson:"name"`
	Country         string               `formType:"country" bson:"country"`
	Code            string               `formType:"string" bson:"code"`
	Type            string               `formType:"string" bson:"type"`
	Description     string               `formType:"string" bson:"description"`
	CompatibleTypes []string             `formType:"list" bson:"compatibles_types"`
	CompatibleWith  []primitive.ObjectID `formType:"compatibility_list" bson:"compatible_with"`
}

// FieldForm holds the defintion of the field on the Form side (UI)
type FieldForm struct {
	Name       string
	Label      string
	Type       string
	Unit       string
	NestedType string
}

// IsNul return true if the the standard is empty
func (s *Standard) IsNul() bool {
	if s.Name != "" && s.Type != "" {
		return false
	}
	return true
}

// GetName returns the Name
func (s *Standard) GetName() string {
	return s.Name
}

func (s *Standard) GetType() string {
	return s.Type
}

// GetCode return the Code
func (s *Standard) GetCode() string {
	return s.Code
}

// GetID return the ID
func (s *Standard) GetID() primitive.ObjectID {
	return s.ID
}

// Delete Standard will remove the Standard struct in the database
func (Standard) Delete(db *mongo.Database, values url.Values, id primitive.ObjectID, standardType StandardInt) (int, interface{}) {
	if id != primitive.NilObjectID {
		log.Print("Will Delete the ID : ")
		log.Println(id)
		collectionName := handledStandard[standardType.GetType()]
		col := db.Collection(collectionName)
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

func (s *Standard) SetID(id primitive.ObjectID) {
	s.ID = id
}

// GetFormFields will return the custom fields of a Standard to be used in the UI
func GetFormFields(s StandardInt) map[string]FieldForm {
	fields := make(map[string]FieldForm)
	t := reflect.TypeOf(s).Elem()
	log.Printf("Type: %s\n", t.Name())
	log.Printf("Kind: %s\n", t.Kind())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// use the json name wichi isn't capitalized
		jsonValue, ok := field.Tag.Lookup("json")
		if ok {
			field.Name = jsonValue
		}
		typeTag := field.Tag.Get("formType")
		if typeTag == "-" {
			continue
		}
		unitTag := field.Tag.Get("formUnit")

		log.Printf("%d, %v (%v), tag : '%v'\n", i, field.Name, field.Type.Name(), typeTag)
		nestedType := ""
		// non standard type
		if typeTag == "nested" || typeTag == "nestedArray" {
			nestedType = field.Type.Name()
			if nestedType == "" {
				nestedType = field.Type.Elem().Name()
			}
		}

		if field.Type.Name() == "" {
			log.Println(field.Type.Elem())
			nestedType = field.Type.Elem().Name()

		}
		fields[field.Name] = FieldForm{
			Name:       field.Name,
			Type:       typeTag,
			Unit:       unitTag,
			NestedType: nestedType,
		}

	}
	return fields
}

// Get Standard return the requests Standards (given the type of standard requested)
func (Standard) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, standardType StandardInt, adj string) (int, interface{}) {
	log.Printf("having Get for standard [%#v] with ID : %d", standardType, id)
	page, err := strconv.ParseInt(values.Get("page"), 10, 64)
	if err != nil {
		page = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	perPage, err := strconv.ParseInt(values.Get("per_page"), 10, 64)
	if err != nil {
		perPage = defaultPerPage
	}

	structOnly := values.Get("struct_only")
	// List possible Compatible types
	compatibleTypes := values.Get("compatible_types_only")
	if compatibleTypes != "" {
		log.Printf("We have a compatible types list request")
		return 200, standardType.GetCompatibleTypes()
	}
	// if the request just want the struct we'll respond an new struct only
	if structOnly != "" {
		log.Print("We have a struct only request")
		return 200, GetFormFields(standardType)
	}
	// Let Display All that We Have
	// Someday Pagination will be there
	if id == primitive.NilObjectID {
		log.Debug("returning every items")
		return 200, GetAll(db, page, perPage, standardType)
	}
	// retrieve collections from handled standards
	collectionName := handledStandard[standardType.GetType()]
	collection := db.Collection(collectionName)
	result := collection.FindOne(context.TODO(), bson.M{"_id": id})

	if result.Err() != nil {
		log.Println(result.Err().Error())
		log.Println(standardType)
		return 404, fmt.Sprintf("Standard not found because %s", result.Err().Error())
	}
	result.Decode(&standardType)

	return 200, standardType
}

// GetAll will returns al the standards
func GetAll(db *mongo.Database, page int64, perPage int64, standardType StandardInt) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	sType := reflect.TypeOf(standardType)
	standards := reflect.New(reflect.SliceOf(sType)).Interface()

	collectionName := handledStandard[standardType.GetType()]
	collection := db.Collection(collectionName)
	skip := page * perPage
	findOpt := options.FindOptions{
		Limit: &perPage,
		Skip:  &skip,
	}
	cursor, err := collection.Find(ctx, bson.M{}, &findOpt)
	if err != nil {
		log.Printf("Whow would not get all standards %s", err.Error())
	}
	cursor.All(ctx, standards)

	return standards
}

// Post will save the Standard
func (Standard) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string, standardType StandardInt) (int, interface{}) {
	body := request.Body
	log.Debugf("Received args : \n\t %+v", values)
	decoder := json.NewDecoder(body)
	standard := reflect.New(reflect.TypeOf(standardType).Elem()).Interface()

	err := decoder.Decode(standard)
	if err != nil {
		log.Printf("Could not unmarshal object : %s", err)
		return http.StatusBadRequest, fmt.Sprintf("Could not unmarshal object : %s", err)
	}

	standardTyped := standard.(StandardInt)
	if standardTyped.IsNul() {
		return http.StatusBadRequest, "The object is null"
	}

	log.Debugf("Saving : %#v", standardTyped)

	standardTyped, err = save(db, standardTyped)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("Could not save the standard : \n\t %s", err.Error())
	}

	return http.StatusCreated, standardTyped
}

// Put updates Standard
func (Standard) Put(db *mongo.Database, values url.Values, body io.ReadCloser, standardType StandardInt) (int, interface{}) {
	log.Debugf("Received args : \n\t %+v", values)
	decoder := json.NewDecoder(body)
	standard := reflect.New(reflect.TypeOf(standardType)).Elem().Interface().(StandardInt)
	err := decoder.Decode(&standard)
	if err != nil {
		log.Printf("Could not unmarshal object : %s", err)
		return http.StatusBadRequest, fmt.Sprintf("Could not unmarshal object : %s", err)

	}
	standardTyped := standard.(StandardInt)
	standardTyped, err = save(db, standardTyped)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("Could not save the standard : \n\t %s", err.Error())
	}
	return http.StatusOK, standardTyped
}

func save(db *mongo.Database, st StandardInt) (StandardInt, error) {
	upsert := false
	if st.GetID() == primitive.NilObjectID {
		st.Init()
		upsert = true
	}
	col := db.Collection(handledStandard[st.GetType()])
	opts := options.ReplaceOptions{
		Upsert: &upsert,
	}
	filter := bson.M{"_id": st.GetID()}
	_, err := col.ReplaceOne(context.TODO(), &filter, st, &opts)
	if err != nil {
		log.Warnf("Could not save the Standard : %s", err.Error())
	}
	return st, err

}
