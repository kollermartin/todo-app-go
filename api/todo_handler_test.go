package api

import (
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
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	db  *sql.DB
	log *logrus.Logger
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

func setupRouter() *gin.Engine {
	router := gin.Default()
	log = logrus.New()

	router.GET("/todos", GetTodos(db, log))

	return router
}

func TestGetTodos(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, 200, w.Code)

	var todos []types.Todo
	err := json.Unmarshal(w.Body.Bytes(), &todos)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	assert.Equal(t, len(testData), len(todos))
}

func TestMain(m *testing.M) {
	testDB := utils.CreateTestDB(testData)

	db = testDB.DbInstance

	defer testDB.CleanUp()

	os.Exit(m.Run())
}
