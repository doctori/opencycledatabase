package main 

import (
	"fmt"
)

func main() {
	var api = new(API)
	bike := new(Bike)
	api.AddResource(bike,"/bike")
	fmt.Printf("Listening To :8080")
	api.Start(8080)
	
}