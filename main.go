package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"urlShorten/controller"
	"urlShorten/database"
	"urlShorten/model"

	"math/rand"

	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

var mainDb *gorm.DB

type ResponseJson struct {
	ShortURL string `json:"short_url"`
}

func main() {

	rand.Seed(time.Now().UnixNano())

	mainDb = database.InitDatabase()

	mainDb.AutoMigrate(&model.URL{})

	r := echo.New()
	r.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			value, _ := c.Get("customValueFromContext").(int)
			fmt.Printf("REQUEST: uri: %v, status: %v, custom-value: %v , METHOD:%v\n", v.URI, v.Status, value, v.Method)
			return nil
		},
	}))
	r.Use(middleware.CORS())
	r.Use(middleware.Recover())

	r.GET("/all-routes", func(c echo.Context) error {
		time.Sleep(2 * time.Second)
		return controller.AllRoutes(c)
	})
	r.GET("/all", func(c echo.Context) error {
		time.Sleep(2 * time.Second)
		return controller.GetAllUrl(c, mainDb)
	})
	r.POST("/shorten", func(c echo.Context) error {
		time.Sleep(2 * time.Second)
		return controller.ShortenURL(c, mainDb)
	})

	r.GET("/:shortURL", func(c echo.Context) error {
		time.Sleep(2 * time.Second)
		return controller.RedirectURL(c, mainDb)
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := r.Start(":8080"); err != nil && err != http.ErrServerClosed {
			r.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Shutdown(ctx); err != nil {
		r.Logger.Fatal(err)
	}
}
