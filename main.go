package main

import (
	db "apitest/database"
	"apitest/handlers"
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"log"
	"net/http"
)

func test(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func main() {
	gin.SetMode(gin.DebugMode)
	database, dbErr := sql.Open("sqlite3", "./test.db")
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	db.InitDb(database)
	router := gin.Default()

	router.GET("/", test)
	router.GET("/users/", func(c *gin.Context) {
		handlers.GetUsers(c, database)
	})
	router.GET("/users/:id", func(c *gin.Context) {
		handlers.UserByID(c, database)
	})
	router.POST("/users/add", func(c *gin.Context) {
		handlers.CreateUser(c, database)
	})
	router.PATCH("/users/update", func(c *gin.Context) {
		handlers.UpdateUser(c, database)
	})

	router.Run("127.0.0.1:8080")
}
