package standards

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

// Crank will define the Crank (not the crankset) standard
type Crank struct {
	Standard `gorm:"embedded" formType:"-"`
	// Length of the crank (cm)
	Length float32 `formType:"int" formUnit:"cm"`
	// BB holds the BB compatibles with this crank
	BB []BottomBracket `gorm:"many2many:cranks_bottombrackets" formType:"nestedArray"`
	// Chainrings will hold the number of compatible chainring standards
	Chainrings []ChainRing `gorm:"many2many:cranks_chainrings" formType:"nestedArray"`
}

func NewCrank() *Crank {
	c := new(Crank)
	c.Type = "Crank"
	return c
}

// Get Crank return the requests Crank Standards ID
func (c *Crank) Get(db *gorm.DB, values url.Values, id int, adj string) (int, interface{}) {
	return c.Standard.Get(db, values, id, c, adj)
}

// Delete Crank delete the requested Crank standard ID
func (c *Crank) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return c.Standard.Delete(db, values, id, c)
}

// Post Crank delete the requested Crank standard ID
func (c *Crank) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return c.Standard.Post(db, values, request, id, adj, c)
}

// Put Crank delete the requested Crank standard ID
func (c *Crank) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return c.Standard.Put(db, values, body, c)
}

// Save Crank will register the crank into the database
func (c *Crank) Save(db *gorm.DB) (err error) {

	// If we have a new record we create it
	if c.GetID() == 0 {
		oldc := new(Crank)
		if errors.Is(db.Where("name = ? AND code = ?", c.GetName(), c.GetCode()).First(&oldc).Error, gorm.ErrRecordNotFound) {
			log.Println("==========================================================================================")
			log.Println("Creating the record")
			log.Printf("%#v", c)
			log.Println("==========================================================================================")

			// We update our just created object in order to add it's associations ...
			err = db.Save(c).Error
			if err != nil {
				return
			}
		} else {
			log.Println("Updating the record")
			log.Println(oldc)
			err = db.Model(oldc).Updates(c).Error
			if err != nil {
				return
			}
			log.Println(oldc)
			// TODO fix this
			db.Model(oldc).First(c, oldc.ID)
		}
	} else {
		err = db.Save(c).Error
		if err != nil {
			return
		}
	}
	return
}
