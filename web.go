package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type DBConfig struct {
	Host     string `json:host`
	Username string `json:username`
	Password string `json:password`
	Port     int    `json:port`
	DBname   string `json:dbname`
}

type Config struct {
	DB DBConfig `json:db`
}

func main() {

	// Load Application Configuration info our Config Struct
	file, err := os.Open("conf/mainConfig.json")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	log.Printf("%#v", config)
	if err != nil {
		log.Fatal(err)
	}
	// Initiate the database =
	db = initDB(config)
	api := new(API)
	bike := new(Bike)
	component := new(Component)
	standard := new(Standard)
	image := new(Image)
	api.AddResource(bike, "/bikes")
	api.AddResource(component, "/components")
	api.AddResource(standard, "/standards")
	api.AddResource(&Brand{}, "/brands")
	api.AddNonJSONResource(image, "")
	fmt.Printf("Listening To :8080 \n")
	api.Start("0.0.0.0", 8080)

}
