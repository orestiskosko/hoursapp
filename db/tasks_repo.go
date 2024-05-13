package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/orestiskosko/hours-app/models"
)

func GetTodaysTaskEntries(conn *pgx.Conn) ([]models.TaskEntry, error) {
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
			WHERE created_at >= date_trunc('day', now())
			AND created_at < date_trunc('day', now()) + interval '1 day'
			`)
	if err != nil {
		return make([]models.TaskEntry, 0), err
	}

	entries, err := pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (models.TaskEntry, error) {
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
		return make([]models.TaskEntry, 0), err
	}

	return entries, nil
}

func GetRunningTaskEntry(conn *pgx.Conn) (models.TaskEntry, error) {
	row := conn.QueryRow(
		context.Background(),
		`
			SELECT * 
			FROM task_entries 
			INNER JOIN tasks ON task_entries.task_id = tasks.id
			INNER JOIN projects ON tasks.project_id = projects.id
			WHERE finished_at IS NULL 
			ORDER BY task_entries.created_at DESC
			LIMIT 1
			`)

	runningTaskEntry := models.TaskEntry{}

	err := row.Scan(
		&runningTaskEntry.ID,
		&runningTaskEntry.TaskID,
		&runningTaskEntry.StartedAt,
		&runningTaskEntry.FinishedAt,
		&runningTaskEntry.Task.ID,
		&runningTaskEntry.Task.Name,
		&runningTaskEntry.Task.Details,
		&runningTaskEntry.Task.Duration,
		&runningTaskEntry.Task.ProjectID,
		&runningTaskEntry.Task.Project.ID,
		&runningTaskEntry.Task.Project.Name,
		&runningTaskEntry.Task.Project.Description,
		&runningTaskEntry.Task.Project.UnitCharge)

	if err != nil && err != pgx.ErrNoRows {
		return runningTaskEntry, err
	}

	return runningTaskEntry, nil
}
