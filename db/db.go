package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetConnection() (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
}

func EnsureMigrated() {
	conn, err := GetConnection()
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	slog.Info("Connected to database: ", "db", conn.Config().Database)

	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			unit_charge INT
		);

		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			details TEXT,
			duration INT NOT NULL DEFAULT 0,
			project_id INT REFERENCES projects
		);

		CREATE TABLE IF NOT EXISTS task_entries (
			id SERIAL PRIMARY KEY,
			task_id INT REFERENCES tasks,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			finished_at TIMESTAMP
		);
	`)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	slog.Info("Created DB tables..")
}
