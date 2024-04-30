package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/orestiskosko/hours-app/db"
	"github.com/orestiskosko/hours-app/models"
	"github.com/orestiskosko/hours-app/templates"
	"github.com/orestiskosko/hours-app/utilities"
)

func UseTasks(e *echo.Echo) {
	g := e.Group("/tasks")

	g.GET("", func(c echo.Context) error {
		conn, err := db.GetConnection()
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		rows, err := conn.Query(
			context.Background(),
			`
			SELECT tasks.id, tasks.name, details, duration, project_id, projects.name as ProjectName FROM tasks
			INNER JOIN projects ON tasks.project_id = projects.id
			`)
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Task])
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		return utilities.Render(c, http.StatusOK, templates.TasksPage(tasks))
	})

	g.GET("/create", func(c echo.Context) error {
		conn, err := db.GetConnection()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		rows, err := conn.Query(context.Background(), `SELECT * FROM projects`)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		projects, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Project])
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return utilities.Render(c, http.StatusOK, templates.CreateTask(projects))
	})

	g.POST("", func(c echo.Context) error {
		name := c.FormValue("name")
		details := c.FormValue("details")
		projectID, err := strconv.Atoi(c.FormValue("project_id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		conn, err := db.GetConnection()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		_, err = conn.Exec(context.Background(),
			`
			INSERT INTO tasks (name, details, project_id) VALUES ($1, $2, $3)
			`,
			name, details, projectID)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		c.Response().Header().Set("HX-Location", `{"path":"/tasks", "target":"#router"}`)
		return c.NoContent(http.StatusOK)
	})
}
