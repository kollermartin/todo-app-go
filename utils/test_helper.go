package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"
	"todo-app/types"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type TestDB struct {
	DbInstance *sql.DB
	Container  testcontainers.Container
}

func SeedDB(db *sql.DB, testData []types.Todo) error {
	for _, todo := range testData {
		_, err := db.Exec("INSERT INTO todos (id, title, created_at) VALUES ($1, $2, $3)", todo.ID, todo.Title, todo.CreatedAt)

		if err != nil {
			return err
		}
	}

	return nil
}

func CreateTestDB(testData []types.Todo) *TestDB {
	var env = map[string]string{
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_DB":       "postgres",
	}

	ctx := context.Background()
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

	if err != nil {
		log.Fatal("Error: Could not start postgres container")
		panic(err)
	}

	port, err := container.MappedPort(ctx, "5432")

	if err != nil {
		log.Fatal("Error: Could not get mapped port")
		panic(err)
	}

	host, err := container.Host(ctx)

	if err != nil {
		log.Fatal("Error: Could not get host")
		panic(err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=postgres password=postgres dbname=postgres sslmode=disable", host, port.Int())

	time.Sleep(time.Second)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(
			"Error: The data source arguments are not valid",
		)
		panic(err)
	}
	time.Sleep(time.Second)

	if err = runMigrations(db, "migrations"); err != nil {
		log.Fatal("Error: Could not run migrations")
		panic(err)
	}

	if (testData != nil) && (len(testData) > 0) {
		time.Sleep(time.Second)

		if err := SeedDB(db, testData); err != nil {
			log.Fatal("Error: Could not seed database")
			panic(err)
		}
	}

	return &TestDB{
		DbInstance: db,
		Container:  container,
	}
}

func (t *TestDB) CleanUp() {
	t.DbInstance.Close()
	if err := t.Container.Terminate(context.Background()); err != nil {
		log.Fatal("Error: Could not terminate container")
		panic(err)
	}
}

func runMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get path")
	}

	pathToMigrationFiles := filepath.Join(filepath.Dir(path), "..", "migrations")

	if err != nil {
		log.Fatalf("Failed to create new database driver instance: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:%s", pathToMigrationFiles),
		"postgres", driver,
	)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
