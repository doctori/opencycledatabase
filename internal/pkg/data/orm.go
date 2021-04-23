package data

import (
	"fmt"
	"log"

	"github.com/doctori/opencycledatabase/internal/pkg/config"
	"github.com/doctori/opencycledatabase/internal/pkg/data/standards"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB will initialise the database connection and the scheme
func InitDB(config *config.Config) *gorm.DB {
	connectionString := fmt.Sprintf(
		"user=%s password='%s' host=%s dbname=%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.DBname)
	log.Printf("Connecting to %s", connectionString)
	db, err := gorm.Open(
		postgres.Open(connectionString),
		&gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	checkErr(err, "Postgres Opening Failed")
	// Debug Mode
	db.Debug()
	db.Logger.LogMode(logger.LogLevel(7))
	db.AutoMigrate(&Image{}, &Brand{}, &Component{}, &Bike{})
	// Standards
	// FIXME : manage to use "managedresources" from api routes
	db.AutoMigrate(
		&standards.BottomBracket{},
		&standards.ChainRing{},
		&standards.Crank{},
		&standards.FrontDerailleur{},
		&standards.RearDerailleur{},
		&standards.Headset{},
		&standards.Hub{},
	)
	checkErr(err, "Create tables failed")

	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}
