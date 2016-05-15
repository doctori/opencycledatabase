package main

import (
	"fmt"
	"github.com/doctori/OpenBicycleDatabase/models"
)

func main() {

	api := new(API)
	bike := new(models.Bike)
	component := new(models.Component)
	standard := new(models.Standard)
	image := new(models.Image)
	api.AddResource(bike, "/bikes")
	api.AddResource(component, "/components")
	api.AddResource(standard, "/standards")
	api.AddResource(&models.Brand{}, "/brands")
	api.AddNonJSONResource(image, "/images")
	fmt.Println("Launching brandToComponent")
	fmt.Printf("Listening To :8080 \n")
	go brandToComponent()
	api.Start("0.0.0.0", 8080)

}
