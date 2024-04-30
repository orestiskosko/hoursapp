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
			return c.NoContent(http.StatusInternalServerError)
		}

		defer conn.Close(context.Background())

		rows, err := conn.Query(context.Background(), "SELECT * FROM projects")
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		defer rows.Close()

		projects, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Project])
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		return utilities.Render(c, http.StatusOK, templates.ProjectsPage(projects))
	})

	g.GET("/:id", func(c echo.Context) error {
		conn, err := db.GetConnection()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		defer conn.Close(context.Background())

		// TODO Do this using join
		rows, err := conn.Query(context.Background(),
			`
			SELECT *
			FROM projects
			WHERE id = $1
			`,
			c.Param("id"))
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		project, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Project])
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		rows, err = conn.Query(context.Background(),
			`
			SELECT *
			FROM tasks
			WHERE project_id = $1
			`,
			c.Param("id"))
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		tasks, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Task])
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
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
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		_, err = conn.Exec(context.Background(), "INSERT INTO projects (name, description, unit_charge) VALUES ($1, $2, $3)", name, description, unitCharge)
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		c.Response().Header().Set("HX-Location", `{"path":"/projects", "target":"#router"}`)
		return c.NoContent(http.StatusOK)
	})

	g.DELETE("/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusBadRequest)
		}

		conn, err := db.GetConnection()
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		_, err = conn.Exec(context.Background(), "DELETE FROM projects WHERE id = $1", id)
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusOK)
	})
}
