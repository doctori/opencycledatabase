package data

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/doctori/opencycledatabase/internal/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitDB will initialise the database connection and the scheme
func InitDB(config *config.Config) *mongo.Database {
	connectionString := fmt.Sprintf("mongodb://%s:%d", config.DB.Host, config.DB.Port)
	log.Infof("Connecting to %s", connectionString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	checkErr(err, "Mongo Opening Failed")
	db := client.Database(config.DB.DBname)
	err = db.Client().Ping(ctx, db.ReadPreference())
	checkErr(err, "Mongo Opening Failed")
	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}
