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

func UseProjects(e *echo.Echo) {
	g := e.Group("/projects")

	g.GET("", func(c echo.Context) error {
		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		defer conn.Close(context.Background())

		rows, err := conn.Query(context.Background(), "SELECT * FROM projects")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer rows.Close()

		projects, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Project])
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return utilities.Render(c, http.StatusOK, templates.ProjectsPage(projects))
	})

	g.GET("/:id", func(c echo.Context) error {
		projectID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		defer conn.Close(context.Background())

		// TODO Do this using join
		rows, err := conn.Query(context.Background(),
			`
			SELECT *
			FROM projects
			WHERE id = $1
			`,
			projectID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		project, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Project])
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		rows, err = conn.Query(context.Background(),
			`
			SELECT id, name, details, duration
			FROM tasks
			WHERE project_id = $1
			`,
			projectID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		tasks := make([]models.Task, 0)

		for rows.Next() {
			currTask := models.Task{}

			err = rows.Scan(&currTask.ID, &currTask.Name, &currTask.Details, &currTask.Duration)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			tasks = append(tasks, currTask)
		}

		return utilities.Render(c, http.StatusOK, templates.ProjectPage(project, tasks))
	})

	g.GET("/create", func(c echo.Context) error {
		return utilities.Render(c, http.StatusOK, templates.CreateProject())
	})

	g.POST("", func(c echo.Context) error {
		name := c.FormValue("name")
		description := c.FormValue("description")
		unitCharge, err := strconv.Atoi(c.FormValue("unit_charge"))
		if err != nil {
			unitCharge = 0
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		_, err = conn.Exec(context.Background(), "INSERT INTO projects (name, description, unit_charge) VALUES ($1, $2, $3)", name, description, unitCharge)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		c.Response().Header().Set("HX-Location", `{"path":"/projects", "target":"#router"}`)
		return c.NoContent(http.StatusOK)
	})

	g.DELETE("/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		_, err = conn.Exec(context.Background(), "DELETE FROM projects WHERE id = $1", id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	})
}
