package standards

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

// FrontDerailleur, holds the specs for a Front Derailleur
type FrontDerailleur struct {
	Standard `gorm:"embedded" formType:"-"`
	// TODO Specs ??
	// BrazedOn : is the Front derailleur need a brazed on fixture or a collar (Clamp On)?
	BrazedOn bool `formType:"bool"`
	// CollarSize
	CollarSize float32 `formType:"int" formUnit:"mm"`
}

func NewFrontDerailleur() *FrontDerailleur {
	fd := new(FrontDerailleur)
	fd.Type = "FrontDerailleur"
	return fd
}

// Get FrontDerailleur return the requests FrontDerailleur Standards ID
func (rd *FrontDerailleur) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return rd.Standard.Get(db, values, id, rd)
}

// Delete FrontDerailleur delete the requested FrontDerailleur standard ID
func (rd *FrontDerailleur) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return rd.Standard.Delete(db, values, id, rd)
}

// Post FrontDerailleur delete the requested FrontDerailleur standard ID
func (rd *FrontDerailleur) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return rd.Standard.Post(db, values, request, id, adj, rd)
}

// Put FrontDerailleur delete the requested FrontDerailleur standard ID
func (rd *FrontDerailleur) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return rd.Standard.Put(db, values, body, rd)
}

// Save FrontDerailleur will register the crank into the database
func (rd *FrontDerailleur) Save(db *gorm.DB) (err error) {

	// If we have a new record we create it
	if rd.GetID() == 0 {
		oldc := new(FrontDerailleur)
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
