package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ladiesman2127/birthdays/internal/app/controllers"
)

func initAuthRouts(userController *controllers.UsersController, api *gin.RouterGroup) {
	api.POST("/signup", userController.SignUp)
	api.POST("/auth", userController.Auth)
}
