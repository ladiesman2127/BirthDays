package app

import (
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ladiesman2127/birthdays/api/server/router"
	"github.com/ladiesman2127/birthdays/internal/app/controllers"
	"github.com/ladiesman2127/birthdays/internal/database"
)

func New() *gin.Engine {
	DB_NAME := os.Getenv("DB_NAME")

	app := gin.Default()

	db := database.New(&DB_NAME)

	userController := controllers.NewUsersController(db)

	router.InitRoutes(userController, app)

	return app
}
