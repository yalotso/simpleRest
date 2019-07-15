package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"
	"simpleRest/apis"
	"simpleRest/app"
	"simpleRest/daos"
	"simpleRest/services"
)

func main() {
	if err := app.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}

	db, err := gorm.Open("mysql", app.Config.DSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.LogMode(false)

	infoLog := log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime)
	logger := app.NewLogger(infoLog, errorLog)

	router := buildRouter(db, logger)
	router.Run(app.Config.ServerPort)
}

func buildRouter(db *gorm.DB, logger app.Logger) *gin.Engine {
	router := gin.New()
	router.Use(
		gin.Recovery(),
		app.OnRequest(logger),
		//app.AccessHandler(log.Printf),
		//app.ErrorHandler(log.Printf),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.Abort()
		c.JSON(http.StatusOK, "OK"+app.Version)
	})
	router.NoRoute(func(c *gin.Context) {
		c.Abort()
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "page not found"})
	})

	router.Use(
		app.Transactional(db),
	)
	rg := router.Group("/v1")
	userDAO := daos.NewUserDAO()
	apis.ServeUserResource(rg, services.NewUserService(userDAO))

	return router
}
