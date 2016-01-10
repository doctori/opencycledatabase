package main

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
)
type Bike struct {
	Name string
	Brand string
	Year int
	Components []Component
	PostNotSupported
    PutNotSupported
    DeleteNotSupported
}
type Component struct {
	Name string
	Brand string
	Type string
}
func (b *Bike) save() error {
	filename := b.Name + ".txt"
	details := b.Name + ";" + b.Brand + ";" + string(b.Year)
	return ioutil.WriteFile(filename, []byte(details), 0600)
}
func (Bike) Get(values url.Values) (int, interface{}) {
	filename := "db/" + values.Get("name") + ".json"
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return 500, ""
	}
	var bike Bike
	if( json.Unmarshal(jsonBlob, &bike)	!= nil){
		return 500, ""
	} 
	return 200,bike
}