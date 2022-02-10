package main

import (
	"github.com/gin-gonic/gin"
	"main.go/Controller"
	"main.go/Models"
)

func main() {	
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	Models.ConnectDataBase()

	//Basic Auth
	/*
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin":"password",
	}))
	*/

	authorized := r.Group("/")
	authorized.Use(Controller.Login)
	{
		authorized.GET("/books", Controller.FindBooks)
	}
	
	r.POST("/books", Controller.CreateBook)
	//r.POST("/login",Controller.Login)
	r.POST("/books/:id", Controller.FindBook)
	r.PATCH("/books/:id", Controller.UpdateBook)
	r.DELETE("/books/:id", Controller.DeleteBook)
	
	r.Run()
}
