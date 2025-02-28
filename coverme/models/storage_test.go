package models

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddTodo(t *testing.T) {
	storage := NewInMemoryStorage()
	todo, err := storage.AddTodo("Test", "Content")
	assert.NoError(t, err)
	assert.Equal(t, "Test", todo.Title)
	assert.Equal(t, "Content", todo.Content)
	assert.False(t, todo.Finished)
}

func TestGetTodo(t *testing.T) {
	storage := NewInMemoryStorage()
	todo, err := storage.AddTodo("Test", "Content")
	assert.NoError(t, err)
	gotTodo, err := storage.GetTodo(todo.ID)
	assert.NoError(t, err)
	assert.Equal(t, todo.ID, gotTodo.ID)
	assert.Equal(t, "Test", gotTodo.Title)
}

func TestGetTodoNotFound(t *testing.T) {
	storage := NewInMemoryStorage()
	_, err := storage.GetTodo(999)
	assert.Error(t, err)
}

func TestFinishTodo(t *testing.T) {
	storage := NewInMemoryStorage()
	todo, err := storage.AddTodo("Test", "Content")
	assert.NoError(t, err)
	err = storage.FinishTodo(todo.ID)
	assert.NoError(t, err)
	finishedTodo, err := storage.GetTodo(todo.ID)
	assert.NoError(t, err)
	assert.True(t, finishedTodo.Finished)
}

func TestFinishTodoNotFound(t *testing.T) {
	storage := NewInMemoryStorage()
	err := storage.FinishTodo(999)
	assert.Error(t, err)
}

func TestGetAllTodos(t *testing.T) {
	storage := NewInMemoryStorage()
	_, err := storage.AddTodo("Title 1", "Content 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = storage.AddTodo("Title 2", "Content 2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	todos, err := storage.GetAll()
	require.NoError(t, err)
	require.Len(t, todos, 2)
}
