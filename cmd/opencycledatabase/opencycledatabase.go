package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/doctori/opencycledatabase/internal/pkg/config"
	"github.com/doctori/opencycledatabase/internal/pkg/data"
	"github.com/doctori/opencycledatabase/internal/pkg/api"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func main() {

	// Initiate the database =

	api := new(api.API)
	bike := new(data.Bike)
	component := new(data.Component)
	bbStandard := new(data.BBStandard)
	image := new(data.Image)
	api.AddResource(bike, "/bikes")
	api.AddResource(component, "/components")
	api.AddResource(bbStandard, "/standards")
	api.AddResource(&data.Brand{}, "/brands")
	api.AddNonJSONResource(image, "")
	fmt.Printf("Listening To :8080 \n")
	api.Start("0.0.0.0", 8080)

}

func init() {
	
	// Load Application Configuration info our Config Struct
	file, err := os.Open("conf/mainConfig.json")
	if err != nil {
		log.Fatal(err)
	}
	var conf config.Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	log.Printf("%#v", conf)
	if err != nil {
		log.Fatal(err)
	}

	db = data.InitDB(conf)
}
