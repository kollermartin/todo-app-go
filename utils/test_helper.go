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

	"github.com/docker/go-connections/nat"
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
	if testData == nil && len(testData) == 0 {
		return nil
	}

	for _, todo := range testData {
		_, err := db.Exec("INSERT INTO todos (external_id, title, created_at) VALUES ($1, $2, $3)", todo.ExternalID, todo.Title, todo.CreatedAt)

		if err != nil {
			return err
		}
	}

	return nil
}

func CreateTestDB(testData []types.Todo) (*TestDB, error) {

	ctx := context.Background()

	container, host, port, err := createTestContainer(ctx)

	if err != nil {
		return nil, err
	}

	db, err := setupDBConnection(host, port)

	if err != nil {
		return nil, err
	}

	if err = runMigrations(db, "migrations"); err != nil {
		return nil, err
	}

	if err := SeedDB(db, testData); err != nil {
		return nil, err
	}

	return &TestDB{
		DbInstance: db,
		Container:  container,
	}, nil
}

func (t *TestDB) CleanUp() {
	t.DbInstance.Close()
	if err := t.Container.Terminate(context.Background()); err != nil {
		log.Fatal("Error: Could not terminate container")
		panic(err)
	}
}

func setupDBConnection(host string, port string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=postgres sslmode=disable", host, port)

	db, err := sql.Open("postgres", connStr)

	// TODO - Fix this
	time.Sleep(time.Second)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTestContainer(ctx context.Context) (container testcontainers.Container, host string, port string, error error) {
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
