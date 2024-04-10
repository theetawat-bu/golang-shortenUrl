package controller

import (
	"fmt"
	"net/http"
	"urlShorten/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RedirectURL(c echo.Context, mainDb *gorm.DB) error {

	shortURL := c.Param("shortURL")
	var myUrl model.URL

	tx := mainDb.Where("url_short = ?", shortURL).Find(&myUrl)
	if tx.Error != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	if tx.RowsAffected == 0 {
		res := fmt.Sprintf("URL short: http://localhost:8080/%v is not found", shortURL)
		return c.String(http.StatusNotFound, res)

	}

	return c.Redirect(http.StatusMovedPermanently, myUrl.UrlFull)

}
