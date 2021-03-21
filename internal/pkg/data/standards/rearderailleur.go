package standards

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

type RearDerailleur struct {
	Standard `gorm:"embedded" formType:"-"`
	// TODO Specs ??
	// CageLength hold the length of the cage in mm
	CageLength float32 `json:"CageLength" formType:"int"  formUnit:"mm"`
}

func NewRearDerailleur() *RearDerailleur {
	rd := new(RearDerailleur)
	rd.Type = "RearDerailleur"
	return rd
}

// Get RearDerailleur return the requests RearDerailleur Standards ID
func (rd *RearDerailleur) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return rd.Standard.Get(db, values, id, rd)
}

// Delete RearDerailleur delete the requested RearDerailleur standard ID
func (rd *RearDerailleur) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return rd.Standard.Delete(db, values, id, rd)
}

// Post RearDerailleur delete the requested RearDerailleur standard ID
func (rd *RearDerailleur) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return rd.Standard.Post(db, values, request, id, adj, rd)
}

// Put RearDerailleur delete the requested RearDerailleur standard ID
func (rd *RearDerailleur) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return rd.Standard.Put(db, values, body, rd)
}

// Save RearDerailleur will register the crank into the database
func (rd *RearDerailleur) Save(db *gorm.DB) (err error) {

	// If we have a new record we create it
	if rd.GetID() == 0 {
		oldc := new(RearDerailleur)
		if errors.Is(db.Where("name = ? AND code = ?", rd.GetName(), rd.GetCode()).First(&oldc).Error, gorm.ErrRecordNotFound) {
			log.Println("==========================================================================================")
			log.Println("Creating the record")
			log.Printf("%#v", rd)
			log.Println("==========================================================================================")

			// We update our just created object in order to add it's associations ...
			err = db.Save(rd).Error
			if err != nil {
				return
			}
		} else {
			log.Println("Updating the record")
			log.Println(oldc)
			err = db.Model(oldc).Updates(rd).Error
			if err != nil {
				return
			}
			log.Println(oldc)
			// TODO fix this
			db.Model(oldc).First(rd, oldc.ID)
		}
	} else {
		err = db.Save(rd).Error
		if err != nil {
			return
		}
	}
	return
}
