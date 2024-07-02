package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ladiesman2127/birthdays/internal/app/controllers"
)

func initUserRoutes(userController *controllers.UsersController, api *gin.RouterGroup) {
	api.GET("user/:id", userController.GetUser)
	api.GET("/users", userController.GetUsers)
	api.GET("/birthdays", userController.GetBirthdayNotifications)
	api.POST("/addFriend/:id", userController.AddFriend)
	api.POST("/removeFriend/:id", userController.RemoveFriend)
}
