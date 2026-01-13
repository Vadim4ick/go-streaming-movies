package routes

import (
	"myapp/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine) {
  router.GET("/movies", controllers.GetMovies())

  router.POST("/register", controllers.RegisterUser())
  router.POST("/login", controllers.LoginUser())
}