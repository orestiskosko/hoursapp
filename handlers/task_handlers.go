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
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		rows, err := conn.Query(
			context.Background(),
			`
			SELECT tasks.id, tasks.name, details, duration, projects.id, projects.name FROM tasks
			INNER JOIN projects ON tasks.project_id = projects.id
			`)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		tasks := make([]models.Task, 0)
		for rows.Next() {
			task := models.Task{}

			err = rows.Scan(&task.ID, &task.Name, &task.Details, &task.Duration, &task.Project.ID, &task.Project.Name)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			tasks = append(tasks, task)
		}

		return utilities.Render(c, http.StatusOK, templates.TasksPage(tasks))
	})

	g.GET("/create", func(c echo.Context) error {
		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		rows, err := conn.Query(context.Background(), `SELECT * FROM projects`)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		projects, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Project])
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return utilities.Render(c, http.StatusOK, templates.CreateTask(projects))
	})

	g.GET("/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())
		task := models.Task{}
		row := conn.QueryRow(
			context.Background(),
			`
			SELECT tasks.id, tasks.name, details, duration, projects.id, projects.name FROM tasks
			INNER JOIN projects ON tasks.project_id = projects.id
			WHERE tasks.id = $1
			`,
			id)

		err = row.Scan(&task.ID, &task.Name, &task.Details, &task.Duration, &task.Project.ID, &task.Project.Name)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return utilities.Render(c, http.StatusOK, templates.EditTask(task))
	})

	g.POST("", func(c echo.Context) error {
		name := c.FormValue("name")
		details := c.FormValue("details")
		projectID, err := strconv.Atoi(c.FormValue("project_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		_, err = conn.Exec(context.Background(),
			`
			INSERT INTO tasks (name, details, project_id) VALUES ($1, $2, $3)
			`,
			name, details, projectID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		c.Response().Header().Set("HX-Location", `{"path":"/tasks", "target":"#router"}`)
		return c.NoContent(http.StatusOK)
	})

	g.PUT("/:id", func(c echo.Context) error {
		name := c.FormValue("name")
		details := c.FormValue("details")
		taskID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		_, err = conn.Exec(context.Background(),
			`
			UPDATE tasks SET name = $1, details = $2 WHERE id = $3
			`,
			name, details, taskID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		c.Response().Header().Set("HX-Location", `{"path":"/tasks", "target":"#router"}`)
		return c.NoContent(http.StatusOK)
	})

	g.DELETE("/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid task ID")
		}

		db, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		_, err = db.Exec(context.Background(), `DELETE FROM tasks WHERE id = $1`, id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.NoContent(http.StatusOK)
	})
}
