package data

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/doctori/opencycledatabase/internal/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitDB will initialise the database connection and the scheme
func InitDB(config *config.Config) *mongo.Database {
	connectionString := fmt.Sprintf("mongodb://127.0.0.1:27017")
	log.Printf("Connecting to %s", connectionString)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	checkErr(err, "Mongo Opening Failed")
	db := client.Database("ocd")
	err = db.Client().Ping(ctx, db.ReadPreference())
	checkErr(err, "Mongo Opening Failed")
	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}
