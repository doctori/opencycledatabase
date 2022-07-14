package auth

import (
	"context"
	"encoding/json"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/doctori/opencycledatabase/internal/pkg/api/utils"
	"github.com/doctori/opencycledatabase/internal/pkg/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	oauth2 "golang.org/x/oauth2"
)

// Login
type Login struct {
	stravaConfig oauth2.Config
	utils.PutNotSupported
	utils.PostNotSupported
	utils.DeleteNotSupported
}

type StravaTokenInfo struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Id        uint32 `json:"id"`
	Profile   string `json:"profile"`
}

var stravaConfig oauth2.Config

//InitStravaConfig will create the Oauth2 configuration for the starvai Oauth2 client
func InitStravaConfig(conf config.AuthConfig) {
	stravaConfig = oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Scopes:       []string{"read"},
		RedirectURL:  conf.RedirectURL,
		// This points to our Authorization Server
		// if our Client ID and Client Secret are valid
		// it will attempt to authorize our user
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.strava.com/oauth/authorize",
			TokenURL: "https://www.strava.com/oauth/token",
		},
	}
}

func (login *Login) Get(db *mongo.Database, values url.Values, id primitive.ObjectID, adj string) (int, interface{}) {
	login.stravaConfig = stravaConfig
	log.Printf("Get Login called\n")
	code := values.Get("code")
	scope := values.Get("scope")
	if code == "" {
		log.Warn("Oauth Code is empty")
		return 400, "Oauth Code is empty"
	}
	if scope == "" {
		log.Warn("Scope is empty")
		return 400, "Scope is empty"
	}
	token, err := login.stravaConfig.Exchange(context.Background(), code)
	if err != nil {
		return 400, err
	}
	//log.Debugf("%#v\n", token)
	data, err := json.Marshal(token.Extra("athlete"))
	if err != nil {
		return 500, err
	}
	var athlete StravaTokenInfo
	err = json.Unmarshal(data, &athlete)
	if err != nil {
		return 500, err
	}
	log.Debugf("L'athlete en question est %s %s with ID %d\n", athlete.Firstname, athlete.Lastname, athlete.Id)
	if err != nil {
		return 400, err
	}
	return 200, token
}
