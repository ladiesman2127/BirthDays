package main

import (
	"fmt"
	"os"

	"github.com/ladiesman2127/birthdays/internal/app"
)

func main() {
	PORT := os.Getenv("PORT")
	app := app.New()
	app.Run(fmt.Sprintf(":%s", PORT))
}
