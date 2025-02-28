package client

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gitlab.com/slon/shad-go/coverme/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Add(t *testing.T) {
	todo := &models.AddRequest{Title: "Test", Content: "Content"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo/create", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req models.AddRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, todo.Title, req.Title)

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&models.Todo{ID: 1, Title: req.Title, Content: req.Content, Finished: false}); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := New(server.URL)
	createdTodo, err := client.Add(todo)
	assert.NoError(t, err)
	assert.Equal(t, "Test", createdTodo.Title)
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo/1", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		todo := &models.Todo{ID: 1, Title: "Test", Content: "Content", Finished: false}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := New(server.URL)
	todo, err := client.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, "Test", todo.Title)
}

func TestClient_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		todos := []*models.Todo{
			{ID: 1, Title: "Test", Content: "Content", Finished: false},
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todos); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	client := New(server.URL)
	todos, err := client.List()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(todos))
}

func TestClient_Finish(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo/1/finish", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL)
	err := client.Finish(1)
	assert.NoError(t, err)
}

func TestClient_Add_Error(t *testing.T) {
	todo := &models.AddRequest{Title: "Test", Content: "Content"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo/create", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := New(server.URL)
	createdTodo, err := client.Add(todo)
	assert.Error(t, err)
	assert.Nil(t, createdTodo)
}

func TestClient_Get_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo/999", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := New(server.URL)
	todo, err := client.Get(999)
	assert.Error(t, err)
	assert.Nil(t, todo)
}

func TestClient_List_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := New(server.URL)
	todos, err := client.List()
	assert.Error(t, err)
	assert.Nil(t, todos)
}

func TestClient_Finish_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo/999/finish", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := New(server.URL)
	err := client.Finish(999)
	assert.Error(t, err)
}

func TestClient_Add_InvalidData(t *testing.T) {
	todo := &models.AddRequest{Title: "", Content: ""}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/todo/create", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := New(server.URL)
	createdTodo, err := client.Add(todo)
	assert.Error(t, err)
	assert.Nil(t, createdTodo)
}
