package tests

import (
	"context"
	"database/sql"
	"log"
	"time"
	"todo-app/internal/core/domain"

	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func SeedDB(db *sql.DB, testData []domain.Todo) error {
	if testData == nil && len(testData) == 0 {
		return nil
	}

	for _, todo := range testData {
		_, err := db.Exec("INSERT INTO todos (title) VALUES ($1)", todo.Title)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetTodosFromDB(db *sql.DB) []domain.Todo {
	var todos []domain.Todo
	rows, err := db.Query("SELECT uuid, id, title, created_at, updated_at FROM todos")

	if err != nil {
		logrus.Error("Error: ", err)
		return nil
	}

	for rows.Next() {
		var todo domain.Todo
		err := rows.Scan(&todo.UUID, &todo.ID, &todo.Title, &todo.CreatedAt, &todo.UpdatedAt)

		if err != nil {
			logrus.Error("Error: ", err)
			return nil
		}

		todos = append(todos, todo)
	}

	return todos

}

func CleanUpContainer(container testcontainers.Container) {
	if err := container.Terminate(context.Background()); err != nil {
		log.Fatal("Error: Could not terminate container")
		panic(err)
	}
}

func CreateTestContainer(ctx context.Context) (container testcontainers.Container, host string, port string, error error) {
	var env = map[string]string{
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_DB":       "postgres",
	}

	var natPort nat.Port

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env:          env,
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	time.Sleep(time.Second)

	if err != nil {
		return nil, "", "", err
	}

	host, err = container.Host(ctx)

	if err != nil {
		return nil, "", "", err
	}

	natPort, err = container.MappedPort(ctx, "5432")

	if err != nil {
		return nil, "", "", err
	}

	port = natPort.Port()

	return container, host, port, nil
}
