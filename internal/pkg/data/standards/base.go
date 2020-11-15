package standards

import "gorm.io/gorm"

// TODO remove this shit from here
// this should be controlled by the API not the datamodel
const defaultPerPage int = 30

// StandardInt interface define all the method that a standard need to have to be a
// real standard struct
type StandardInt interface {
	GetName() string
	GetCountry() string
	GetCode() string
	GetType() string
	Get() string
}

// Standard define the generic common Standard properties
type Standard struct {
	// add basic ID/Created@/Updated@/Delete@ through Gorm
	gorm.Model
	// TODO : this is a embded struct on other structs,
	//  find a way to create indexes on each struct that embed this struct
	Name        string `gorm:"uniqueIndex:standard_uniqueness"`
	Country     string
	Code        string `gorm:"uniqueIndex:standard_uniqueness"`
	Type        string `gorm:"uniqueIndex:standard_uniqueness"`
	Description string
}

// ThreadStandard defines the the standard of a thread used
// todo : https://en.wikipedia.org/wiki/Screw_thread
type ThreadStandard struct {
	gorm.Model
	ThreadPerInch float32
	Diameter      float32
	Orientation   string
}

func (s *Standard) isNul() bool {
	if s.Name != "" && s.Type != "" {
		return false
	}
	return true
}
