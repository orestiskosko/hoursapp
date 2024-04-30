package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"

	"github.com/orestiskosko/hours-app/db"
	"github.com/orestiskosko/hours-app/models"
	"github.com/orestiskosko/hours-app/templates"
	"github.com/orestiskosko/hours-app/utilities"
)

func UseTracker(e *echo.Echo) {
	g := e.Group("/tracker")

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

		return utilities.Render(c, http.StatusOK, templates.TrackerPage(projects))
	})

	g.GET("/tasks-select", func(c echo.Context) error {
		conn, err := db.GetConnection()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		defer conn.Close(context.Background())

		rows, err := conn.Query(context.Background(), "SELECT * FROM tasks WHERE project_id = $1", c.QueryParam("project_id"))
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		defer rows.Close()

		tasks, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Task])
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusInternalServerError)
		}

		return utilities.Render(c, http.StatusOK, templates.TaskOptions(tasks))
	})

	g.POST("/start", func(c echo.Context) error {
		taskId, err := strconv.Atoi(c.FormValue("task_id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		taskEntry := models.TaskEntry{
			TaskID:    taskId,
			StartedAt: time.Now().UTC(),
		}

		err = conn.QueryRow(context.Background(),
			`INSERT INTO task_entries (task_id, created_at) VALUES ($1, $2) RETURNING id`,
			taskEntry.TaskID,
			taskEntry.StartedAt).Scan(&taskEntry.ID)
		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusBadRequest)
		}

		return utilities.Render(
			c,
			http.StatusOK,
			templates.Timer(
				strconv.Itoa(taskEntry.ID),
				true,
				taskEntry.StartedAt.Format(time.RFC3339)))
	})

	g.POST("/stop", func(c echo.Context) error {
		taskEntryId, err := strconv.Atoi(c.FormValue("task_entry_id"))
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		_, err = conn.Exec(context.Background(),
			`UPDATE task_entries SET finished_at = $1 WHERE id = $2`,
			time.Now().UTC(),
			taskEntryId)

		if err != nil {
			c.Logger().Error(err.Error())
			return c.NoContent(http.StatusBadRequest)
		}

		return utilities.Render(c, http.StatusOK, templates.Timer(strconv.Itoa(taskEntryId), false, ""))
	})

	g.POST("/track", func(c echo.Context) error {
		taskEntryId := c.FormValue("task_entry_id")
		slog.Info(taskEntryId)
		return c.NoContent(http.StatusOK)
	})
}
