package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/orestiskosko/hours-app/db"
	"github.com/orestiskosko/hours-app/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	db.EnsureMigrated()

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Static("/", "public")

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err.Error())
		event := fmt.Sprintf(`{"showToaster": { "message": "%s" } }`, err.Error())
		c.Response().Header().Set("HX-Trigger", event)
		c.Response().Header().Set("HX-Reswap", "none")
		c.Response().Status = http.StatusNoContent
	}

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/tracker")
	})

	handlers.UseTracker(e)
	handlers.UseProjects(e)
	handlers.UseTasks(e)

	e.Logger.Fatal(e.Start(":3000"))
}
