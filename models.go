package main

import (
	"encoding/json"
	"io/ioutil"
)
type Bike struct {
	Name string
	Brand string
	Year int
	Components []Component
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
func bikeLoad(name string) (*Bike, error) {
	filename := "db/" + name + ".json"
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var bike Bike
	err = json.Unmarshal(jsonBlob, &bike)
	
	return &bike, err
}