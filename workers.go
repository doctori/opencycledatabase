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
	for i := 0; i < totalComponents; i = i + offset {
		fmt.Printf("Counting To : %d\n", i)
		//time.Sleep(time.Duration(1) * time.Second)
		models.DB.Limit(offset).Offset(i).Find(&components)
		//Basic string contained test ... to be improved ? surely
		go componentMatchBrands(&brands, &components)

	}
	fmt.Println("brandToComponent Ended")
}

func componentMatchBrands(brands *[]models.Brand, components *[]models.Component) {
	// We have a slice of components and a slice of Brand, let's find if we have a matching pair
	// If we have let's update our component
	for _, component := range *components {
	BRANDS:
		for _, brand := range *brands {
			//Should Ignore Case
			if strings.Contains(component.Name, brand.Name) {
				fmt.Printf("FOUND Component : %s, for the Brand %s\n", component.Name, brand.Name)
				go extractBrandFromComponent(&brand, &component)
				break BRANDS
			}
		}
	}
}

// Wil extract the Brand from the Component Name and save that component
// This might be SQL intensive ?
func extractBrandFromComponent(brand *models.Brand, component *models.Component) {
	fmt.Printf("Updating the Component %s\n", component.Name)
	index := strings.Index(component.Name, brand.Name)
	if index != -1 {
		component.Name = strings.Replace(component.Name, brand.Name, "", 1)
		if len(component.Name) > 4 {
			component.Name = strings.TrimSpace(component.Name)
			component.Brand = *brand
			component.Save()
		}

	}
}
