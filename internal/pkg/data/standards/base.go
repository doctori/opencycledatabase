package standards

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"gorm.io/gorm"
)

// TODO remove this shit from here
// this should be controlled by the API not the datamodel
const defaultPerPage int = 30

// StandardInt interface define all the method that a standard need to have to be a
// real standard struct
type StandardInt interface {
	GetName() string
	//	GetCountry() string
	GetCode() string
	GetID() uint
	//	GetType() string
	//	Get() string
	IsNul() bool
	Get(db *gorm.DB, values url.Values, id int) (int, interface{})
	Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{})
	Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{})
	Delete(db *gorm.DB, values url.Values, id int) (int, interface{})
	Save(db *gorm.DB) (err error)
}

// Standard define the generic common Standard properties
type Standard struct {
	// add basic ID/Created@/Updated@/Delete@ through Gorm
	gorm.Model
	// TODO : this is a embded struct on other structs,
	//  find a way to create indexes on each struct that embed this struct
	Name        string
	Country     string
	Code        string
	Type        string
	Description string
}

// ThreadStandard defines the the standard of a thread used
// todo : https://en.wikipedia.org/wiki/Screw_thread
type ThreadStandard struct {
	gorm.Model
	ThreadPerInch float32
	Diameter      float32
	Orientation   string
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

// GetCode return the Code
func (s *Standard) GetCode() string {
	return s.Code
}

// GetID return the ID
func (s *Standard) GetID() uint {
	return s.ID
}

// Delete Standard will remove the Standard struct in the database
func (Standard) Delete(db *gorm.DB, values url.Values, id int, standardType StandardInt) (int, interface{}) {
	// retrieve the Standard
	standard := reflect.New(reflect.TypeOf(standardType))
	err := db.First(&standard, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Standard not found"
	}
	// TODO : implement Delete method

	db.Delete(standard)
	return 200, ""
}

// Get Standard return the requests Standards (given the type of standard requested)
func (Standard) Get(db *gorm.DB, values url.Values, id int, standardType StandardInt) (int, interface{}) {
	log.Printf("having Get for standard [%#v] with ID : %d", standardType, id)
	page := values.Get("page")
	perPage := values.Get("per_page")

	//standard := reflect.New(reflect.TypeOf(standardType)).Interface()
	// Let Display All that We Have
	// Someday Pagination will be there
	if id == 0 {
		return 200, GetAll(db, page, perPage, standardType)
	}
	err := db.Model(standardType).First(&standardType, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Standard not found"
	}
	return 200, standardType
}

// GetAll will returns al the standards
func GetAll(db *gorm.DB, page string, perPage string, standardType StandardInt) interface{} {
	ipage, err := strconv.Atoi(page)
	if err != nil {
		ipage = 0
	}
	// Retrieve the per_page arg, if not a number default to 30
	iperPage, err := strconv.Atoi(perPage)
	if err != nil {
		iperPage = defaultPerPage
	}

	//db.Preload("Components").Preload("Components.Standards").Find(&bikes) // Don't Need to load every Component for the main List
	sType := reflect.TypeOf(standardType)
	standards := reflect.New(reflect.SliceOf(sType)).Interface()

	//var cr []ChainRing
	log.Printf("%#v\n", standards)
	db.Model(standardType).
		Offset(ipage * iperPage).
		Limit(iperPage).
		Find(standards)
	log.Printf("%#v\n", standards)
	return standards
}

// Post will save the BBStandard
func (Standard) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string, standardType StandardInt) (int, interface{}) {
	body := request.Body
	fmt.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	log.Println(reflect.TypeOf(standardType))
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
	err = standardTyped.Save(db)
	if err != nil {
		return http.StatusInternalServerError, fmt.Sprintf("Could Not Save the Standard : \n\t %s", err.Error())
	}
	return http.StatusAccepted, standardTyped
}

// Put updates Standard
func (Standard) Put(db *gorm.DB, values url.Values, body io.ReadCloser, standardType StandardInt) (int, interface{}) {
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
