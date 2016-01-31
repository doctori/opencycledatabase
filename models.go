package main

import (
	"encoding/json"
	"io/ioutil"
	"io"
	"log"
	"net/url"
	"fmt"
)
type Bike struct {
	Name string
	Brand string
	Year int
	Components []Component
    PutNotSupported
    DeleteNotSupported
}
type Component struct {
	Name string
	Brand string
	Type string
	Standard string
	Year int
	PostNotSupported
	PutNotSupported
	DeleteNotSupported
}


func (b Bike) save() (error) {
	filename := "db/"+b.Name + ".json"
	jsonBlob, err := json.Marshal(b)
	if (err != nil ){
		return err
	}
	return ioutil.WriteFile(filename, jsonBlob, 0600)
}
func (Bike) Get(values url.Values) (int, interface{}) {
	filename := "db/" + values.Get("name") + ".json"
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
		fmt.Printf("Received args : \n\t %+v\n", values)
		decoder := json.NewDecoder(body)
    	var bike Bike
    	err := decoder.Decode(&bike)
    	if err != nil {
        	panic("shiit")
    	}
    	log.Println(bike)
		err = bike.save()
		if (err != nil){
			return 500, "Could Not Save the Bike"
		}
	   	return 200,bike

	}
func (Bike) addComponent(bike Bike, component Component) {

}
func (Component) Get(values url.Values) (int, interface{}) {
	filename := "db/"+values.Get("name") + ".json"
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