package main

import (
	"fmt"
	"myapp/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()

  r.GET("/movies", controllers.GetMovies())
  r.GET("/movies/:movie_id", controllers.GetMovie())
  r.POST("/movies", controllers.AddMovie())
  
  r.POST("/register", controllers.RegisterUser())
  
  if err := r.Run(":8080"); err != nil {
	fmt.Println("Failed to start server", err)
  }
}