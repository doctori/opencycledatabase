package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/doctori/opencycledatabase/internal/pkg/api"
	"github.com/doctori/opencycledatabase/internal/pkg/config"
	"github.com/doctori/opencycledatabase/internal/pkg/data"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db   *mongo.Database
	conf *config.Config
)

func main() {

	// Initiate the database =

	api := new(api.API)
	api.Init(db, conf)

}

func init() {

	// Load Application Configuration info our Config Struct
	file, err := os.Open("conf/mainConfig.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	log.Printf("%#v", conf)
	if err != nil {
		log.Fatal(err)
	}

	db = data.InitDB(conf)
}
