package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ladiesman2127/birthdays/internal/app/controllers"
)

func InitRoutes(userController *controllers.UsersController, app *gin.Engine) {
	api := app.Group("api")
	initAuthRouts(userController, api)
	initUserRoutes(userController, api)
}
