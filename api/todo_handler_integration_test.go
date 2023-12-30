package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"todo-app/types"
	"todo-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	db       *sql.DB
	log      *logrus.Logger
	router   *gin.Engine
	testData = []types.Todo{
		{ID: "2233a6b2-ae99-40fc-bdd7-db49834993ab", Title: "Task 1", CreatedAt: time.Date(2023, 12, 29, 18, 26, 45, 0, time.UTC)},
		{ID: "1c15f5f7-3207-4d4a-b50f-f6f8bacfb0e9", Title: "Task 2", CreatedAt: time.Date(2023, 12, 29, 18, 25, 18, 0, time.UTC)},
		{ID: "4ffaaf6e-6693-45a4-b1d2-02da81bebc46", Title: "Task 3", CreatedAt: time.Date(2023, 12, 29, 18, 32, 19, 0, time.UTC)},
		{ID: "6e6a468c-bb56-4d11-ab57-b91c46501ae7", Title: "Task 4", CreatedAt: time.Date(2023, 12, 29, 18, 16, 29, 0, time.UTC)},
		{ID: "1c5b7f6f-ee0b-4e48-b246-c8206d1dccc2", Title: "Task 5", CreatedAt: time.Date(2023, 12, 29, 18, 13, 57, 0, time.UTC)},
		{ID: "0c3cd173-ec71-42f8-a191-49bf5613f3f0", Title: "Task 6", CreatedAt: time.Date(2023, 12, 29, 18, 37, 1, 0, time.UTC)},
		{ID: "71757bd0-21fe-44f0-8768-4be10fd2e8e5", Title: "Task 7", CreatedAt: time.Date(2023, 12, 29, 18, 30, 40, 0, time.UTC)},
		{ID: "2ea5cc31-fe70-444c-887d-b48a22d8f265", Title: "Task 8", CreatedAt: time.Date(2023, 12, 29, 18, 29, 22, 0, time.UTC)},
		{ID: "597b2371-bd2a-48cc-8c25-e018a37803f4", Title: "Task 9", CreatedAt: time.Date(2023, 12, 29, 18, 14, 28, 0, time.UTC)},
		{ID: "400e27ed-32ff-4e3a-b6e7-0e0c09a0c121", Title: "Task 10", CreatedAt: time.Date(2023, 12, 29, 18, 6, 58, 0, time.UTC)},
	}
)

func setupRouter() {
	router = gin.Default()
	log = logrus.New()

	router.GET("/todos", GetTodos(db, log))
	router.GET("/todos/:id", GetTodoByID(db, log))
	router.POST("/todos", CreateTodo(db, log))
	router.PUT("/todos/:id", UpdateTodo(db, log))

}

func TestGetTodos(t *testing.T) {

	req, _ := http.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Run("It should get all todos", func(t *testing.T) {
		var todos []types.Todo
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
		todo := testData[0]

		req, _ := http.NewRequest("GET", "/todos/"+todo.ID, nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var todoResponse types.Todo

		err := json.Unmarshal(w.Body.Bytes(), &todoResponse)

		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.EqualValues(t, todo, todoResponse)
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

		todoInput := types.TodoInput{
			Title: "Test todo",
		}

		jsonValue, _ := json.Marshal(todoInput)

		req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonValue))

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)

		var todoResponse types.Todo

		err := json.Unmarshal(w.Body.Bytes(), &todoResponse)

		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.Equal(t, todoInput.Title, todoResponse.Title)
		assert.NotEmpty(t, todoResponse.ID)
		assert.NotEmpty(t, todoResponse.CreatedAt)
	})

	t.Run("It should return 400 if title is missing", func(t *testing.T) {
		todoInput := types.TodoInput{
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
		todo := testData[0]

		todoInput := types.TodoInput{
			Title: "Updated task",
		}

		jsonValue, _ := json.Marshal(todoInput)

		req, _ := http.NewRequest("PUT", "/todos/"+todo.ID, bytes.NewBuffer(jsonValue))

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var todoResponse types.Todo

		err := json.Unmarshal(w.Body.Bytes(), &todoResponse)

		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.EqualValues(t, todo.ID, todoResponse.ID)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("It should return 404 if todo doesnt exist", func(t *testing.T) {
		randomUUID := uuid.New().String()
		todoInput := types.TodoInput{
			Title: "Updated task",
		}

		jsonValue, _ := json.Marshal(todoInput)

		req, _ := http.NewRequest("PUT", "/todos/"+randomUUID, bytes.NewBuffer(jsonValue))

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestMain(m *testing.M) {
	testDB, error := utils.CreateTestDB(testData)

	if error != nil {
		panic(error.Error())
	}

	db = testDB.DbInstance

	setupRouter()

	defer testDB.CleanUp()

	os.Exit(m.Run())
}
