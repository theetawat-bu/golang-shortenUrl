package controller

import (
	"net/http"
	"urlShorten/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetAllUrl(c echo.Context, mainDb *gorm.DB) error {
	lock.Lock()
	defer lock.Unlock()
	var allUrl []model.URL

	if tx := mainDb.Find(&allUrl); tx.Error != nil {
		return c.String(http.StatusBadRequest, "Error Find All URl")
	}

	return c.JSON(http.StatusOK, allUrl)

}
