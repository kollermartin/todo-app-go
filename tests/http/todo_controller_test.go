package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	// "time"
	"todo-app/config"
	handler "todo-app/internal/adapter/handler/http"
	"todo-app/internal/adapter/postgres"
	"todo-app/internal/adapter/postgres/repository"
	"todo-app/internal/core/domain"
	"todo-app/internal/core/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

var (
	db       *sql.DB
	router   *gin.Engine
	testData = []domain.Todo{
		{Title: "Task 1"},
		{Title: "Task 2"},
		{Title: "Task 3"},
		{Title: "Task 4"},
		{Title: "Task 5"},
		{Title: "Task 6"},
		{Title: "Task 7"},
		{Title: "Task 8"},
		{Title: "Task 9"},
		{Title: "Task 10"},
	}
)

func Init() (db *postgres.DB, container testcontainers.Container) {
	ctx := context.Background()

	container, host, port, error := CreateTestContainer(ctx)
	if error != nil {
		logrus.Fatal("Error creating test container", error)
	}

	config := config.Config{
		App: &config.App{
			Name:           "todo-app-test",
			Env:            "test",
			Port:           "3000",
			MigrationsPath: "../../migrations",
		},
		Db: &config.Db{
			Host:     host,
			Port:     port,
			User:     "postgres",
			Password: "postgres",
			Name:     "postgres",
			Type:     "postgres",
		},
		HTTP: &config.HTTP{
			URL:  "localhost",
			Port: "8080",
		},
	}

	db, err := postgres.New(ctx, config.Db)
	if err != nil {
		logrus.Fatal("Error initializing test database", err)
	}

	err = db.Migrate(config.App)
	if err != nil {
		logrus.Fatal("Error running migrations", err)
	}

	if err := SeedDB(db.SqlDB, testData); err != nil {
		logrus.Fatal("Error seeding database", err)
	}

	return db, container
}

func TestGetTodos(t *testing.T) {

	req, _ := http.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Run("It should get all todos", func(t *testing.T) {
		var todos []handler.TodoResponse
		err := json.Unmarshal(w.Body.Bytes(), &todos)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, len(testData), len(todos))
	})
}

func TestGetTodoByID(t *testing.T) {
	t.Run("It should return todo by id", func(t *testing.T) {
		todos := GetTodosFromDB(db)

		expectedTodoRes := handler.NewTodoResponse(&todos[0])

		req, _ := http.NewRequest("GET", "/todos/"+expectedTodoRes.ID.String(), nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var todoResponse handler.TodoResponse

		err := json.Unmarshal(w.Body.Bytes(), &todoResponse)

		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedTodoRes.ID, todoResponse.ID)
		assert.Equal(t, expectedTodoRes.Title, todoResponse.Title)
		assert.Equal(t, expectedTodoRes.CreatedAt.In(time.UTC), todoResponse.CreatedAt.In(time.UTC))
		assert.Equal(t, expectedTodoRes.UpdatedAt.In(time.UTC), todoResponse.UpdatedAt.In(time.UTC))
	})

	t.Run("It should return 400 if todo id is not uuid", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/todos/123", nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("It should return 404 if todo is not found", func(t *testing.T) {
		randomUUID := uuid.New().String()
		req, _ := http.NewRequest("GET", "/todos/"+randomUUID, nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCreateTodo(t *testing.T) {
	t.Run("It should create a new todo", func(t *testing.T) {

		todoInput := handler.CreateRequest{
			Title: "Test todo",
		}

		jsonValue, _ := json.Marshal(todoInput)

		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonValue))

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)

		var todoResponse handler.TodoResponse

		err := json.Unmarshal(w.Body.Bytes(), &todoResponse)

		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, todoInput.Title, todoResponse.Title)
		assert.NotEmpty(t, todoResponse.ID)
		assert.NotEmpty(t, todoResponse.CreatedAt)
	})

	t.Run("It should return 400 if title is missing", func(t *testing.T) {
		todoInput := handler.CreateRequest{
			Title: "",
		}

		jsonValue, _ := json.Marshal(todoInput)

		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonValue))

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}

func TestUpdateTodo(t *testing.T) {
	t.Run("It should return 400 if todo id is not uuid", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/todos/123", nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("It should update todo", func(t *testing.T) {
		todos := GetTodosFromDB(db)

		todo := handler.NewTodoResponse(&todos[0])

		todoInput := handler.UpdateRequest{
			Title: "Updated task",
		}

		jsonValue, _ := json.Marshal(todoInput)

		req, _ := http.NewRequest("PUT", "/todos/"+todo.ID.String(), bytes.NewBuffer(jsonValue))

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var todoResponse handler.TodoResponse

		err := json.Unmarshal(w.Body.Bytes(), &todoResponse)

		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, todo.ID, todoResponse.ID)
		assert.Equal(t, todoInput.Title, todoResponse.Title)
		assert.Equal(t, todo.CreatedAt.In(time.UTC), todoResponse.CreatedAt.In(time.UTC))
		assert.NotEqual(t, todo.UpdatedAt.In(time.UTC), todoResponse.UpdatedAt.In(time.UTC))
		assert.True(t, todoResponse.UpdatedAt.After(todo.UpdatedAt))
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("It should return 404 if todo doesnt exist", func(t *testing.T) {
		randomUUID := uuid.New().String()
		todoInput := handler.UpdateRequest{
			Title: "Updated task",
		}

		jsonValue, _ := json.Marshal(todoInput)

		req, _ := http.NewRequest("PUT", "/todos/"+randomUUID, bytes.NewBuffer(jsonValue))

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteTodo(t *testing.T) {
	t.Run("It should return 400 if todo id is not uuid", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/todos/123", nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("It should delete todo", func(t *testing.T) {
		todos := GetTodosFromDB(db)
		todo := handler.NewTodoResponse(&todos[0])

		req, _ := http.NewRequest("DELETE", "/todos/"+todo.ID.String(), nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("It should return 404 if todo doesnt exist", func(t *testing.T) {
		randomUUID := uuid.New().String()

		req, _ := http.NewRequest("DELETE", "/todos/"+randomUUID, nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestMain(m *testing.M) {

	testDb, container := Init()

	db = testDb.SqlDB

	defer testDb.Close()
	defer CleanUpContainer(container)

	todoRepo := repository.NewTodoRepository(testDb)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)
	routeris, err := handler.NewRouter(todoHandler)
	if err != nil {
		logrus.Fatal("Error initializing router", err)
	}

	router = routeris.Engine

	os.Exit(m.Run())
}
