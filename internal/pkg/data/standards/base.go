package standards

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO remove this shit from here
// this should be controlled by the API not the datamodel
const defaultPerPage int = 30

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
	//GetCompatible() []StandardInt
	Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{})
	Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string) (int, interface{})
	Put(db *mongo.Database, values url.Values, body io.ReadCloser) (int, interface{})
	Delete(db *mongo.Database, values url.Values, id primitive.ObjectID) (int, interface{})
	Save(db *mongo.Database) (err error)
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
	page := values.Get("page")
	perPage := values.Get("per_page")

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
		log.Print("returning every items")
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
func GetAll(db *mongo.Database, page string, perPage string, standardType StandardInt) interface{} {
	/*ipage, err := strconv.Atoi(page)
	if err != nil {
		ipage = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	iperPage, err := strconv.Atoi(perPage)
	if err != nil {
		iperPage = defaultPerPage
	}
	*/
	//db.Preload("Components").Preload("Components.Standards").Find(&bikes) // Don't Need to load every Component for the main List
	sType := reflect.TypeOf(standardType)
	standards := reflect.New(reflect.SliceOf(sType)).Interface()

	log.Printf("%#v\n", standards)
	collectionName := handledStandard[standardType.GetType()]
	collection := db.Collection(collectionName)
	// TODO : pagination
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Whow would not get all standards %s", err.Error())
	}
	cursor.All(context.TODO(), standards)
	return standards
}

// Post will save the Standard
func (Standard) Post(db *mongo.Database, values url.Values, request *http.Request, id primitive.ObjectID, adj string, standardType StandardInt) (int, interface{}) {
	body := request.Body
	log.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	log.Println(reflect.TypeOf(standardType))
	standard := reflect.New(reflect.TypeOf(standardType).Elem()).Interface()

	err := decoder.Decode(standard)
	if err != nil {
		log.Printf("Could not unmarshal object : %s", err)
		return http.StatusBadRequest, fmt.Sprintf("Could not unmarshal object : %s", err)
	}

	standardTyped := standard.(StandardInt)
	stdStandard := Standard{
		Type: standardTyped.GetType(),
		Name: standardTyped.GetName(),
	}
	stdStandard.ID = standardTyped.GetID()
	log.Printf("Saving : %#v", standardTyped)
	if standardTyped.IsNul() {
		return http.StatusBadRequest, "The object is null"
	}
	// we need to save the linked between the "standards" table and the typed standard table
	err = standardTyped.Save(db)
	if err != nil {
		log.Printf("Could not save the standard : \n\t %s", err.Error())
		return http.StatusInternalServerError, fmt.Sprintf("Could not save the standard : \n\t %s", err.Error())
	}
	return http.StatusAccepted, standardTyped
}

func (s *Standard) save(db *mongo.Database) (err error) {
	collectionName := handledStandard[s.GetType()]
	col := db.Collection(collectionName)
	if s.ID == primitive.NilObjectID {
		s.ID = primitive.NewObjectID()
		log.Printf("Object of type %s is new inserting it into collection %s", s.GetType(), collectionName)
		var res = &mongo.InsertOneResult{}
		res, err = col.InsertOne(context.TODO(), s)
		log.Print(res)
		return
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": s.ID}
	_, err = col.UpdateOne(context.TODO(), filter, s, opts)
	return
}

// Put updates Standard
func (Standard) Put(db *mongo.Database, values url.Values, body io.ReadCloser, standardType StandardInt) (int, interface{}) {
	fmt.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	standard := reflect.New(reflect.TypeOf(standardType)).Elem().Interface().(StandardInt)
	err := decoder.Decode(&standard)
	if err != nil {
		log.Printf("Could not unmarshal object : %s", err)
		return http.StatusBadRequest, fmt.Sprintf("Could not unmarshal object : %s", err)

	}
	log.Println(standard)
	err = standard.Save(db)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("Could Not Save the Standard : \n\t%s", err.Error())
	}
	return http.StatusOK, standard
}
