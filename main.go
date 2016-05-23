package main

import (
	"flag"
	"fmt"
	"github.com/doctori/OpenBicycleDatabase/models"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "Debug Mode")
	flag.BoolVar(&debug, "d", false, "Debug Mode")

}
func main() {
	flag.Parse()
	if debug {
		fmt.Println("DEBUG MODE")
		models.DebugMode()
	}
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
