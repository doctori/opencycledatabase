package data

import (
	"errors"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Image is the description and the pointer to the image
type Image struct {
	gorm.Model
	Name          string `gorm:"uniqueIndex:image_uniqueness"`
	Path          string `gorm:"uniqueIndex:image_uniqueness"`
	Type          string
	ContentType   string
	ContentLength int64
	Content       []byte `sql:"-"`
	PutNotSupported
	DeleteNotSupported
}

// Get will return the image (content type is image/jpeg)
// TODO : make the content detection working
func (Image) Get(db *gorm.DB, values url.Values, id int) (int, interface{}) {
	var img Image
	if id == 0 {
		return 200, img.List(db)
	}
	err := db.First(&img, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, "Image Not Found"
	}
	file, err := os.Open(img.Path)
	log.Printf("Serving : %v", img.Path)
	defer file.Close()
	if err != nil {
		log.Println("Could Not Find the Image")
		return 404, "File Not Found"
	}
	// Detect Mime Type
	buff := make([]byte, 512)
	_, err = file.ReadAt(buff, 0)
	if err != nil {
		log.Panic(err)
	}
	img.ContentType = http.DetectContentType(buff)
	log.Printf("ContentType is : %s", img.ContentType)
	return 200, img

}

func (i Image) getContentType() string {
	return i.ContentType
}

func (i Image) getContentLength() string {
	return strconv.FormatInt(i.ContentLength, 10)
}

func (i Image) getContent() []byte {
	return i.Content
}

// Post saves the images to the "upload" directory (shouldn't it go to S3 ? )
func (Image) Post(db *gorm.DB, values url.Values, request *http.Request, id int, adj string) (int, interface{}) {
	uploadFolder := "upload/"
	mediaType, params, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	checkErr(err, "Could Not Determine the Content Type")
	img := Image{}
	if strings.HasPrefix(mediaType, "multipart/") {
		log.Println("YOU ARE POSTING AN IMAGE")
		multipartReader := multipart.NewReader(request.Body, params["boundary"])
		// Should Build buffer and Write it at the end (loop thour the nextpart !)
		partReader, err := multipartReader.NextPart()
		fileName := partReader.FileName()
		filePath := uploadFolder + fileName
		file, err := os.Create(filePath)
		if err != nil {
			log.Println(err)
			return 500, err
		}
		checkErr(err, "Could Not Get the Next Part")
		if _, err := io.Copy(file, partReader); err != nil {
			log.Fatal(err)
			return 500, err
		}
		// Let's record this image and return it to our client
		img = Image{Name: fileName, Path: filePath}
		recordedImg := img.save(db)
		log.Printf("Posted : %#v", img)
		return 200, recordedImg
	}

	return 200, img
}
func (Image) List(db *gorm.DB) (images []Image) {
	result := db.Find(&images)
	if result.Error != nil {
		log.Printf("Could not list all images ! because %s", result.Error.Error())
	}
	return
}
func (i Image) save(db *gorm.DB) Image {
	if i.ID == 0 {
		oldi := new(Image)
		db.Where("name = ? and path = ? ", i.Name, i.Path).First(&oldi)
		if oldi.Name == "" {
			log.Println("Recording the New Image")
			db.Create(&i)
		} else {
			log.Println("Updating The Image Record")
			db.Model(&oldi).Updates(&i)
			i = *oldi
			log.Printf("Saving Image %#v", i)
		}
	} else {
		db.Create(&i)
	}
	return i
}
