package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/orestiskosko/hours-app/models"
)

func GetProjects(conn *pgx.Conn) ([]models.Project, error) {
	projects := make([]models.Project, 0)

	rows, err := conn.Query(
		context.Background(),
		`SELECT * FROM projects`)
	if err != nil {
		return projects, err
	}

	for rows.Next() {
		project := models.Project{}
		rows.Scan(&project.ID, &project.Name, &project.Description, &project.UnitCharge)
		projects = append(projects, project)
	}

	return projects, nil
}
