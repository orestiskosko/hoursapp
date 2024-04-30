package models

import "time"

type Project struct {
	ID          int
	Name        string
	Description string
	UnitCharge  int
}

type Task struct {
	ID        int
	Name      string
	Details   string
	Duration  int
	ProjectID int
}

type TaskEntry struct {
	ID         int
	TaskID     int
	StartedAt  time.Time
	FinishedAt *time.Time
}
