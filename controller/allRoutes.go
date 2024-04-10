package controller

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func AllRoutes(e echo.Context) error {
	data, err := json.MarshalIndent(e.Echo().Routes(), "", "  ")
	if err != nil {
		return e.String(http.StatusBadRequest, "")
	}
	os.WriteFile("routes.json", data, 0644)
	return e.JSONBlob(http.StatusOK, data)
}
