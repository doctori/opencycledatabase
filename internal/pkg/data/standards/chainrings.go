package standards

import (
	"errors"
	"io"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

// ChainRing will define the Chainrings standard
type ChainRing struct {
	Standard `gorm:"embedded" formType:"-"`
	// BoltCircleDiameter of the chainring (ref : https://www.sheldonbrown.com/gloss_bo-z.html#bcd)
	BoltCircleDiameter float32 `formType:"int" formUnit:"cm"`
	// BoltsNumber hold the number of bolt on the chainring
	BoltsNumber int `formType:"int" formUnit:"count"`
	// IsIntegrated is true if the chainring is soldered to the crank
	IsIntegrated bool `formType:"bool"`
	// IsDirectlyMounted weither the chainring is a direct mount chainring or not
	IsDirectlyMounted bool `formType:"bool"`
	// Teeth is the number of teeth a chainring has (0-255)
	Teeth uint8 `formType:"int" formUnit:"count"`
}

// NewChainRing return a ChainRing empty object with some predefined fields
func NewChainRing() *ChainRing {
	cr := new(ChainRing)
	cr.Type = "ChainRing"
	return cr
}

// Get ChainRing return the requests ChainRing Standards ID
func (cr *ChainRing) Get(db *gorm.DB, values url.Values, id int, adj string) (int, interface{}) {
	return cr.Standard.Get(db, values, id, cr, adj)
}

// Delete ChainRing delete the requested ChainRing standard ID
func (cr *ChainRing) Delete(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	return cr.Standard.Delete(db, values, id, cr)
}

// Post ChainRing delete the requested ChainRing standard ID
func (cr *ChainRing) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	return cr.Standard.Post(db, values, request, id, adj, cr)
}

// Put ChainRing delete the requested ChainRing standard ID
func (cr *ChainRing) Put(db *gorm.DB, values url.Values, body io.ReadCloser) (int, interface{}) {
	return cr.Standard.Put(db, values, body, cr)
}

// Save ChainRing will register the BB into the database
func (cr *ChainRing) Save(db *gorm.DB) (err error) {

	// If we have a new record we create it
	if cr.GetID() == 0 {
		oldcr := new(ChainRing)
		if errors.Is(db.Where("name = ? AND code = ?", cr.GetName(), cr.GetCode()).First(&oldcr).Error, gorm.ErrRecordNotFound) {
			// We update our just created object in order to add it's associations ...
			err = db.Save(cr).Error
			if err != nil {
				return
			}
		} else {
			err = db.Model(oldcr).Updates(cr).Error
			if err != nil {
				return
			}
			db.Model(oldcr).First(cr, oldcr.ID)
		}
	} else {
		err = db.Save(cr).Error
		if err != nil {
			return
		}
	}
	return
}
