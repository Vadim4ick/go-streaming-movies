package main

import (
	"fmt"
	"myapp/routes"

	"github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()

  routes.SetupUnProtectedRoutes(r)
  routes.SetupProtectedRoutes(r)
  
  if err := r.Run(":8080"); err != nil {
	fmt.Println("Failed to start server", err)
  }
}