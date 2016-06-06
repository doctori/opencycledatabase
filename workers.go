package main

import (
	"fmt"
	"github.com/doctori/OpenBicycleDatabase/models"
	"strings"
)

// Will Parse Comonent names and tries to find Brand Name
// If Found the brand names are extracted from the Component Name and the component is associated with it
// It's based on the brand list whithin the database
func brandToComponent() {
	fmt.Println("brandToComponent Running")
	//TODO
	var brands []models.Brand
	var components []models.Component
	var offset = 5
	models.DB.Find(&brands)
	// Let's cycle through the component 5 by 5
	var totalComponents int
	models.DB.Model(&models.Component{}).Count(&totalComponents)
	fmt.Println(totalComponents)
	componentsChannel := make(chan []models.Component)
	done := make(chan bool)
	// pop a few workers
	for w := 1; w <= 5; w++ {
		go componentMatchBrands(w, componentsChannel, &brands, done)
	}
	for i := 0; i < totalComponents; i = i + offset {
		//time.Sleep(time.Duration(1) * time.Second)
		models.DB.Limit(offset).Offset(i).Find(&components)
		//Basic string contained test ... to be improved ? surely
		componentsChannel <- components
	}

	fmt.Println("brandToComponent Ended")
	for i := 0; i < totalComponents; i = i + offset {
		<-done
	}
}

func componentMatchBrands(workerId int, componentsChannel <-chan []models.Component, brands *[]models.Brand, done chan bool) {
	// We have a slice of components and a slice of Brand, let's find if we have a matching pair
	// If we have let's update our component
	for components := range componentsChannel {
		for _, component := range components {
		BRANDS:
			for _, brand := range *brands {
				//Should Ignore Case
				if brand.Name != "" {
					if strings.Contains(component.Name, brand.Name) {
						fmt.Printf("[WORKER %d] FOUND Component : %s, for the Brand %s\n", workerId, component.Name, brand.Name)
						extractBrandFromComponent(&brand, &component)
						break BRANDS
					}
				}
			}
		}
	}
	done <- true
}

// Wil extract the Brand from the Component Name and save that component
// This might be SQL intensive ?
func extractBrandFromComponent(brand *models.Brand, component *models.Component) {
	fmt.Printf("Updating the Component %s\n", component.Name)
	index := strings.Index(component.Name, brand.Name)
	if index != -1 {
		if len(component.Name) > 4 && len(brand.Name) > 4 {
			component.Name = strings.Replace(component.Name, brand.Name, "", 1)
			component.Name = strings.TrimSpace(component.Name)
			component.Brand = *brand
			component.Save()
		}

	}

}
