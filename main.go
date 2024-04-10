package main

import (
	"fmt"
	"net/http"
	"urlShorten/model"

	"math/rand"

	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var urlMap = make(map[string]string)
var mainDb *gorm.DB

func main() {

	rand.Seed(time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open("url.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	mainDb = db
	// Migrate the schema
	mainDb.AutoMigrate(&model.URL{})
	r := echo.New()

	r.POST("/shorten", shortenURL)

	r.GET("/:shortURL", redirectURL)
	r.GET("/all", allUrl)

	r.Start(":8080")

}

type ResponseJson struct {
	ShortURL string `json:"short_url"`
}

func allUrl(c echo.Context) error {
	var allUrl []model.URL
	tx := mainDb.Find(&allUrl)
	if tx.Error != nil {
		return c.String(http.StatusBadRequest, "Error Find All URl")
	}
	return c.JSON(http.StatusNotFound, allUrl)

}
func checkIfShorten(longURL string) error {
	var myUrl model.URL
	fmt.Println(longURL)

	tx := mainDb.Where("url_full = ?", longURL).Take(&myUrl)
	fmt.Println(myUrl)

	fmt.Println(tx.RowsAffected)
	if tx.RowsAffected != 0 {
		res := fmt.Sprintf("URL %v has been shorten", longURL)
		return errors.New(res)

	}
	return nil
}
func shortenURL(c echo.Context) error {

	longURL := c.FormValue("url")
	haveFound := checkIfShorten(longURL)
	if haveFound != nil {
		return c.String(http.StatusBadRequest, haveFound.Error())
	}

	shortURL := generateShortURL()
	myUrl := model.URL{UrlFull: longURL, UrlShort: shortURL}

	tx := mainDb.Create(&myUrl)
	if tx.Error != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	return c.JSON(http.StatusOK, myUrl)

}

func redirectURL(c echo.Context) error {

	shortURL := c.Param("shortURL")
	var myUrl model.URL

	tx := mainDb.Where("url_short = ?", shortURL).Find(&myUrl)
	if tx.Error != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	if tx.RowsAffected == 0 {
		res := fmt.Sprintf("URL short: http://localhost/%v is not found", shortURL)
		return c.String(http.StatusNotFound, res)

	}

	return c.Redirect(http.StatusMovedPermanently, myUrl.UrlFull)

}

func generateShortURL() string {

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	shortURL := make([]byte, 6)

	for i := range shortURL {

		shortURL[i] = charset[rand.Intn(len(charset))]

	}

	return string(shortURL)

}
