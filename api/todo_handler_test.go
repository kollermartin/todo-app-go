package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"todo-app/test"
	"todo-app/types"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	db  *sql.DB
	log *logrus.Logger
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

	assert.Equal(t, 0, len(todos))
}

func TestMain(m *testing.M) {
	testDB := test.CreateTestDB()

	db = testDB.DbInstance

	defer testDB.CleanUp()

	os.Exit(m.Run())
}
