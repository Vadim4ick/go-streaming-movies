package main

import (
	"fmt"
	"myapp/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()

 // Define a simple GET endpoint
  r.GET("/hello", func(c *gin.Context) {
    // Return JSON response
	c.String(200, "HELLO")
  })

  r.GET("/movies", controllers.GetMovies())
  r.GET("/movies/:movie_id", controllers.GetMovie())
  
  if err := r.Run(":8080"); err != nil {
	fmt.Println("Failed to start server", err)
  }
}