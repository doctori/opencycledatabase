package main 

import (
	"fmt"
)

func main() {
	var api = new(API)
	bike := new(Bike)
	component := new(Component)
	standard := new(Standard)
	api.AddResource(bike,"/bikes")
	api.AddResource(component, "/components")
	api.AddResource(standard, "/standards")
	fmt.Printf("Listening To :8080 \n")
	api.Start("0.0.0.0",8080)
	
}