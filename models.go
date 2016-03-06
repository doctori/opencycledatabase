package main

import (
	"encoding/json"
	"io/ioutil"
	"io"
	"os"
	"log"
	"net/url"
	"strings"
	"fmt"
	"path/filepath"
	"github.com/satori/go.uuid"
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

func (b Bike) save() (error) {
	filename := "db/bike_"+b.Name + ".json"
	jsonBlob, err := json.Marshal(b)
	if (err != nil ){
		return err
	}
	return ioutil.WriteFile(filename, jsonBlob, 0600)
}

func GetAllBikes()(interface{}){
	var bikes []Bike
	files, err := ioutil.ReadDir("db/")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if (strings.Index(file.Name(),"bike_") == 0 ) {
			absPath, err := filepath.Abs("db/"+file.Name())
			jsonBlob, err := ioutil.ReadFile(absPath)
			if err != nil {
				return err
			}
			var bike Bike
			if( json.Unmarshal(jsonBlob, &bike)	!= nil) {
				return ""
			}
			bikes = append(bikes,bike)
			fmt.Println(file.Name())
		}
		
	}
	return bikes
}
func (Bike) Get(values url.Values) (int, interface{}) {
	if values.Get("name") == "" {
		return 200,GetAllBikes();
	}
	filename := "db/bike_" + values.Get("name") + ".json"
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return 404, "404 Bike Not Found"
	}
	var bike Bike
	if( json.Unmarshal(jsonBlob, &bike)	!= nil){
		return 500, ""
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
	files, err := ioutil.ReadDir("db/")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if (strings.Index(file.Name(),"component_") == 0 ) {
			absPath, err := filepath.Abs("db/"+file.Name())
			jsonBlob, err := ioutil.ReadFile(absPath)
			if err != nil {
				panic(err)
			}
			var c Component
			if( json.Unmarshal(jsonBlob, &c)	!= nil) {
				panic("Could not Unmarshal !")
			}
			components = append(components,c)
			fmt.Println(file.Name())
		}		
	}
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
	filename := "db/component_"+c.Name + ".json"
	if (c.Id == ""){
		c.Id = uuid.NewV4().String()
	}
	jsonBlob, err := json.Marshal(c)
	if (err != nil ){
		return err,c
	}
	// Retrieve the Supported Standard ? OR not !
	err = ioutil.WriteFile(filename, jsonBlob, 0600)
	return err,c
}
