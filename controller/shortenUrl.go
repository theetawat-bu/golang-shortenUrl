package controller

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"urlShorten/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	lock = sync.Mutex{}
)

func ShortenURL(c echo.Context, mainDb *gorm.DB) error {
	lock.Lock()
	defer lock.Unlock()
	longURL := c.FormValue("url")
	haveFound := checkIfShorten(longURL, mainDb)
	if haveFound != nil {
		return c.String(http.StatusBadRequest, haveFound.Error())
	}

	shortURL := generateShortURL()
	myUrl := model.URL{UrlFull: longURL, UrlShort: shortURL}
	tx := mainDb.Create(&myUrl)
	if tx.Error != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	return c.JSON(http.StatusCreated, myUrl)

}

func checkIfShorten(longURL string, mainDb *gorm.DB) error {
	var myUrl model.URL

	tx := mainDb.Where("url_full = ?", longURL).Take(&myUrl)

	if tx.RowsAffected != 0 {
		res := fmt.Sprintf("URL %v has been shorten", longURL)
		return errors.New(res)

	}
	return nil
}

func generateShortURL() string {

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	shortURL := make([]byte, 7)

	for i := range shortURL {
		shortURL[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortURL)

}
