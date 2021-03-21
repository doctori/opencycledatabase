package standards

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

type Spoke struct {
	Standard `gorm:"embedded" formType:"-"`
	// Length is the length of the spoke
	Length int16 `formType:"int" formUnit:"mm"`
}

func NewSpoke() *Spoke {
	s := new(Spoke)
	s.Type = "Spoke"
	return s
}

// Get Spoke return the requests Spoke Standards ID
func (s *Spoke) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return s.Standard.Get(db, values, id, s)
}

// Delete Spoke delete the requested Spoke standard ID
func (s *Spoke) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return s.Standard.Delete(db, values, id, s)
}

// Post Spoke delete the requested Spoke standard ID
func (s *Spoke) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return s.Standard.Post(db, values, request, id, adj, s)
}

// Put Spoke delete the requested Spoke standard ID
func (s *Spoke) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return s.Standard.Put(db, values, body, s)
}

// Save Spoke will register the crank into the database
func (s *Spoke) Save(db *gorm.DB) (err error) {

	// If we have a new record we create it
	if s.GetID() == 0 {
		oldc := new(Spoke)
		if errors.Is(db.Where("name = ? AND code = ?", s.GetName(), s.GetCode()).First(&oldc).Error, gorm.ErrRecordNotFound) {
			log.Println("==========================================================================================")
			log.Println("Creating the record")
			log.Printf("%#v", s)
			log.Println("==========================================================================================")

			// We update our just created object in order to add it's associations ...
			err = db.Save(s).Error
			if err != nil {
				return
			}
		} else {
			log.Println("Updating the record")
			log.Println(oldc)
			err = db.Model(oldc).Updates(s).Error
			if err != nil {
				return
			}
			log.Println(oldc)
			// TODO fix this
			db.Model(oldc).First(s, oldc.ID)
		}
	} else {
		err = db.Save(s).Error
		if err != nil {
			return
		}
	}
	return
}
