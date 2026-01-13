package routes

import (
	"myapp/controllers"
	"myapp/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
  router.Use(middleware.AuthMiddleware())

  router.GET("/movies/:movie_id", controllers.GetMovie())
  router.POST("/movies", controllers.AddMovie())
}