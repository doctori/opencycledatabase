package main

import (
	"encoding/json"
	"io/ioutil"
	"io"
	"os"
	"log"
	"net/url"
	"fmt"
	"github.com/satori/go.uuid"
	"labix.org/v2/mgo"
	"gopkg.in/mgo.v2/bson"
)
type Bike struct {
	Name string
	Brand string
	Year int
	Components []Component
	SupportedStandard []Standard
    PutNotSupported
    DeleteNotSupported
}
type Component struct {
	Id string
	Name string
	Brand string
	Type string
	Standards []Standard 
	Year int
	PutNotSupported
	DeleteNotSupported
}
type Standard struct {
	Id string
	Name string
	Country string
	Code string
	Type string
}
var session, _  = mgo.Dial("mongodb://127.0.0.1/OpenBicycleDatabase")
	
func (b Bike) save() (error) {

	bikeCollection := session.DB("OpenBicycleDatabase").C("bikes")
	// Check if element exits
	filter := bson.M{"name":b.Name}
	bike, err := bikeCollection.Upsert(filter,b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received upsert : \n\t %+v\n", bike)
	//bike already exists, let's update it
	return err
}

func GetAllBikes()(interface{}){
	
	var bikes []Bike
	bikeCollection := session.DB("OpenBicycleDatabase").C("bikes")
	bikeCollection.Find(bson.M{}).All(&bikes)
	return bikes
}
func (Bike) Get(values url.Values) (int, interface{}) {
	if values.Get("name") == "" {
		return 200,GetAllBikes();
	}
	bikeCollection := session.DB("OpenBicycleDatabase").C("bikes")
	query := bikeCollection.Find(bson.M{"name":values.Get("name")})
	var bike Bike
	err := query.One(&bike)
	if err != nil {
		if err == mgo.ErrNotFound {
			return 404, "Bike not found"
		}else{
			return 500, "Internal Server Error"
			panic(err)	
		}
		
	}

	if err != nil {
		return 404, "404 Bike Not Found"
	}

	return 200,bike
}

func (Bike) Post(values url.Values, body io.ReadCloser) (int, interface{}) {
		decoder := json.NewDecoder(body)
    	var bike Bike
    	err := decoder.Decode(&bike)
    	if err != nil {
        	panic(err)
        	return 500, "Internal Error"
    	}
    	log.Println(bike)
		err = bike.save()
		if (err != nil){
			return 500, "Could Not Save the Bike"
		}
	   	return 200,bike

	}
func (Bike) addComponent(bike Bike, component Component) {
	bike.Components = append(bike.Components,component)
	bike.save()
}
func (Bike) getCompatibleComponents(bike Bike) ([]Component){
	var compatibleComponents []Component
	// Get The Components, Get their Standards

	// Get All the Components supporting those standards
	// Return Thems
	return compatibleComponents
}

func (Standard) Post(values url.Values, body io.ReadCloser) (int, interface{}) {
	fmt.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	var standard Standard
	err := decoder.Decode(&standard)
	if err != nil {
    	panic("shiit")
	}
	log.Println(standard)
	err,standard = standard.save()
	if (err != nil){
		return 500, "Could Not Save the Standard"
	}
   	return 200,standard
}

func (Standard) Put(values url.Values, body io.ReadCloser) (int, interface{}) {
	fmt.Printf("Received args : \n\t %+v\n", values)
	decoder := json.NewDecoder(body)
	var standard Standard
	err := decoder.Decode(&standard)
	if err != nil {
    	panic("shiit")
	}
	log.Println(standard)
	err,standard = standard.save()
	if (err != nil){
		return 500, "Could Not Save the Standard"
	}
   	return 200,standard
}
func (Standard) Get(values url.Values) (int, interface{}) {
	filename := "db/standard_" + values.Get("name") + ".json"
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return 404, "404 Standard Not Found"
	}
	var s Standard
	if( json.Unmarshal(jsonBlob, &s)	!= nil){
		return 500, ""
	} 

	return 200,s
}
func (Standard) Delete(values url.Values) (int, interface{}) {
	filename := "db/standard_" + values.Get("name") + ".json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return 404, "404 Standard Not Found"
	}
	if( os.Remove(filename)	!= nil){
		return 500, "Could not Delete the Standard"
	}
	return 200,""
}
func (s Standard) save() (error, Standard){
	filename := "db/standard_"+s.Name + ".json"
	if (s.Id == ""){
		s.Id = uuid.NewV4().String()
	}
	jsonBlob, err := json.Marshal(s)
	if (err != nil ){
		return err,s
	}
	err = ioutil.WriteFile(filename, jsonBlob, 0600)
	return err,s
}

func (Component) Post(values url.Values, body io.ReadCloser) (int, interface{}) {
	fmt.Printf("Received args : \n\t %+v\n", body)
	decoder := json.NewDecoder(body)
	var component Component
	err := decoder.Decode(&component)
	if err != nil {
    	panic("shiit")
	}
	log.Println(component)
	err,component = component.save()
	if (err != nil){
		return 500, "Could Not Save the Component"
	}
   	return 200,component
}

func (Component) Get(values url.Values) (int, interface{}) {
	if values.Get("name") == "" {
		c := new (Component)
		return 200,c.getAll();
	}
	filename := "db/component_"+values.Get("name") + ".json"
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return 404, "404 Component not Found"
	}
	var component Component
	if ( json.Unmarshal(jsonBlob, &component )!= nil){
		return 500, ""
	}
	return 200,component
}
func (Component) getAll() ([]Component) {
	var components []Component
	componentCollection := session.DB("OpenBicycleDatabase").C("components")
	componentCollection.Find(bson.M{}).All(&components)
	return components
}
// Return Components Compatibles with a standards
func (Component) getCompatible(s Standard)([]Component){
	var compatibleComponents []Component
	// Crado Way
	c := new (Component)
	components := c.getAll()
	for _, component := range components {
		for _, standard := range component.Standards {
			if standard == s {
				compatibleComponents = append(compatibleComponents, component)
			}
		}
	}
	return compatibleComponents
}

func (c Component) save() (error, Component){

	componentCollection := session.DB("OpenBicycleDatabase").C("components")
	// Check if element exits
	filter := bson.M{"name":c.Name}
	component, err := componentCollection.Upsert(filter,c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received upsert : \n\t %+v\n", component)
	//bike already exists, let's update it
	return err,c
}