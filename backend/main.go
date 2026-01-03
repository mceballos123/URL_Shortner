package main

import (

	"github.com/gin-gonic/gin"
	"github.com/mceballos123/URL_Shortner/backend/db"
	"github.com/mceballos123/URL_Shortner/backend/api"
)

func main(){
	router := gin.New()
	db:= db.ConnectDB()

	router.POST("/create-url",api.PostCreateUrl(db))
	router.POST("/create-user",api.PostCreateUser(db))
	router.POST("/login",api.PostLogin(db))

	router.GET("/urls",api.GetAllUrls(db))
	router.GET("/urls/:id",api.GetUrlsByID(db))
	router.GET("/users",api.GetAllUsers(db))
	router.GET("/users/:id",api.GetUserByID(db))

	router.Run(":8080")
}