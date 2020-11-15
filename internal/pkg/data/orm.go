package data

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/doctori/opencycledatabase/internal/pkg/config"
)

var db = &gorm.DB{}

// InitDB will initialise the database connection and the scheme
func InitDB(config config.Config) *gorm.DB {
	connectionString := fmt.Sprintf(
		"user=%s password='%s' host=%s dbname=%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.DBname)
	log.Printf("Connecting to %s", connectionString)
	db, err := gorm.Open("postgres", connectionString)
	checkErr(err, "Postgres Opening Failed")
	// Debug Mode
	db.LogMode(true)
	db.CreateTable(&Image{}, &ComponentType{}, &Brand{}, &BBStandard{}, &Component{}, &Bike{})
	db.Model(&Bike{}).AddUniqueIndex("bike_uniqueness", "name,  year")
	db.Model(&Component{}).AddUniqueIndex("component_uniqueness", "name, year")
	db.Model(&BBStandard{}).AddUniqueIndex("standard_uniqueness", "name, code, type")
	db.Model(&Image{}).AddUniqueIndex("image_uniqueness", "name", "path")
	db.AutoMigrate(&Bike{}, &Component{}, &BBStandard{}, &Image{}, &Brand{}, &ComponentType{})
	checkErr(err, "Create tables failed")

	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}
