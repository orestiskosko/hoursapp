package handlers

import (
	"context"
	"fmt"
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
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		projects, err := db.GetProjects(conn)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		projectOptions := make(map[string]string)
		for _, project := range projects {
			projectOptions[fmt.Sprintf("%d", project.ID)] = project.Name
		}

		todayEntries, err := db.GetTodaysTaskEntries(conn)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		runningTaskEntry, err := db.GetRunningTaskEntry(conn)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		date := time.Now().UTC().Format(time.RFC3339)
		if runningTaskEntry.ID != 0 {
			date = runningTaskEntry.StartedAt.Format(time.RFC3339)
		}

		return utilities.Render(
			c,
			http.StatusOK,
			templates.TrackerPage(models.TrackerViewModel{
				IsRunning:        runningTaskEntry.ID != 0,
				Date:             date,
				RunningTaskEntry: runningTaskEntry.ToViewModel(),
				TaskEntries:      models.ToViewModels(todayEntries),
				ProjectOptions:   projectOptions,
			}))
	})

	g.GET("/entries", func(c echo.Context) error {
		dateRaw := c.QueryParam("date")
		date, err := time.Parse("2006-01-02", dateRaw)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		fmt.Printf("\n\n%s\n\n", date.Format(time.RFC3339))

		rows, err := conn.Query(
			context.Background(),
			`
			SELECT 
				task_entries.id, 
				task_entries.task_id, 
				task_entries.created_at, 
				task_entries.finished_at, 
				tasks.id,
				tasks.name,
				tasks.details,
				tasks.duration,
				tasks.project_id,
				projects.id,
				projects.name,
				projects.description,
				projects.unit_charge
			FROM task_entries 
			INNER JOIN tasks ON task_entries.task_id = tasks.id
			INNER JOIN projects ON tasks.project_id = projects.id
			WHERE created_at >= date_trunc('day', $1::timestamp)
			AND created_at < date_trunc('day', $1::timestamp) + interval '1 day'
			`,
			date.Format(time.RFC3339))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		todayEntries, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.TaskEntry, error) {
			taskEntry := models.TaskEntry{}

			err := row.Scan(
				&taskEntry.ID,
				&taskEntry.TaskID,
				&taskEntry.StartedAt,
				&taskEntry.FinishedAt,
				&taskEntry.Task.ID,
				&taskEntry.Task.Name,
				&taskEntry.Task.Details,
				&taskEntry.Task.Duration,
				&taskEntry.Task.ProjectID,
				&taskEntry.Task.Project.ID,
				&taskEntry.Task.Project.Name,
				&taskEntry.Task.Project.Description,
				&taskEntry.Task.Project.UnitCharge)

			if err != nil {
				return taskEntry, err
			}

			return taskEntry, nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return utilities.Render(c, http.StatusOK, templates.TaskEntries(models.ToViewModels(todayEntries)))
	})

	g.GET("/tasks-select", func(c echo.Context) error {
		projectID, err := strconv.Atoi(c.QueryParam("project_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		defer conn.Close(context.Background())

		rows, err := conn.Query(
			context.Background(),
			"SELECT id, name, details FROM tasks WHERE project_id = $1",
			projectID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer rows.Close()

		taskOptions := make(map[string]string)
		for rows.Next() {
			var taskID int
			var taskName string
			var taskDetails string
			err = rows.Scan(&taskID, &taskName, &taskDetails)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			taskOptions[fmt.Sprintf("%d", taskID)] = taskName
		}

		return utilities.Render(c, http.StatusOK, templates.TaskOptions(taskOptions))
	})

	g.POST("/start", func(c echo.Context) error {
		taskId, err := strconv.Atoi(c.FormValue("task_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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

		row := conn.QueryRow(
			context.Background(),
			`
			WITH inserted_task AS 
			(
				INSERT INTO task_entries (task_id, created_at) VALUES ($1, $2)
				RETURNING *
			) 
			SELECT *
			FROM inserted_task
			INNER JOIN tasks ON inserted_task.task_id = tasks.id
			INNER JOIN projects ON tasks.project_id = projects.id
			`,
			taskEntry.TaskID,
			taskEntry.StartedAt)

		err = row.Scan(
			&taskEntry.ID,
			&taskEntry.TaskID,
			&taskEntry.StartedAt,
			&taskEntry.FinishedAt,
			&taskEntry.Task.ID,
			&taskEntry.Task.Name,
			&taskEntry.Task.Details,
			&taskEntry.Task.Duration,
			&taskEntry.Task.ProjectID,
			&taskEntry.Task.Project.ID,
			&taskEntry.Task.Project.Name,
			&taskEntry.Task.Project.Description,
			&taskEntry.Task.Project.UnitCharge)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return utilities.Render(
			c,
			http.StatusOK,
			templates.StartTimerResponse(taskEntry.ToViewModel()))
	})

	g.POST("/stop", func(c echo.Context) error {
		taskEntryId, err := strconv.Atoi(c.FormValue("task_entry_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		conn, err := db.GetConnection()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer conn.Close(context.Background())

		row := conn.QueryRow(context.Background(),
			`
			UPDATE task_entries 
			SET finished_at = $1 
			FROM tasks, projects
			WHERE task_entries.id = $2 
			RETURNING *`,
			time.Now().UTC(),
			taskEntryId)
		taskEntry := models.TaskEntry{}
		err = row.Scan(
			&taskEntry.ID,
			&taskEntry.TaskID,
			&taskEntry.StartedAt,
			&taskEntry.FinishedAt,
			&taskEntry.Task.ID,
			&taskEntry.Task.Name,
			&taskEntry.Task.Details,
			&taskEntry.Task.Duration,
			&taskEntry.Task.ProjectID,
			&taskEntry.Task.Project.ID,
			&taskEntry.Task.Project.Name,
			&taskEntry.Task.Project.Description,
			&taskEntry.Task.Project.UnitCharge)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		rows, err := conn.Query(
			context.Background(),
			`
			SELECT id, name FROM projects
			`)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		projectOptions := make(map[string]string)
		for rows.Next() {
			var projectId int
			var projectName string
			rows.Scan(&projectId, &projectName)
			projectOptions[fmt.Sprintf("%d", projectId)] = projectName
		}

		return utilities.Render(c, http.StatusOK, templates.StopTimerResponse(projectOptions, taskEntry.ToViewModel()))
	})
}
