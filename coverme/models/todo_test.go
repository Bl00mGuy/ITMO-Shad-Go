package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTodoEmptyFields(t *testing.T) {
	todo := &Todo{
		ID:       2,
		Title:    "",
		Content:  "",
		Finished: false,
	}

	assert.Empty(t, todo.Title, "Title should be empty")
	assert.Empty(t, todo.Content, "Content should be empty")
	assert.False(t, todo.Finished, "Finished should be false by default")
}

func TestMarkFinished(t *testing.T) {
	todo := &Todo{
		ID:       1,
		Title:    "Test Todo",
		Content:  "This is a test",
		Finished: false,
	}

	todo.MarkFinished()

	assert.True(t, todo.Finished)
}

func TestMarkUnfinished(t *testing.T) {
	todo := &Todo{
		ID:       1,
		Title:    "Test Todo",
		Content:  "This is a test",
		Finished: true,
	}

	todo.MarkUnfinished()

	assert.False(t, todo.Finished)
}

func TestTodoInitialization(t *testing.T) {
	todo := &Todo{
		ID:       1,
		Title:    "Test Todo",
		Content:  "This is a test",
		Finished: false,
	}

	assert.Equal(t, "Test Todo", todo.Title)
	assert.Equal(t, "This is a test", todo.Content)
	assert.False(t, todo.Finished)
}

func TestTodoIDBoundaries(t *testing.T) {
	todoMin := &Todo{
		ID:       ID(-1),
		Title:    "Negative ID",
		Content:  "This is fine",
		Finished: false,
	}

	todoMax := &Todo{
		ID:       ID(2147483647),
		Title:    "Max ID",
		Content:  "This is fine too",
		Finished: false,
	}

	assert.Equal(t, ID(-1), todoMin.ID)
	assert.Equal(t, ID(2147483647), todoMax.ID)
}

func TestTodoStatusToggle(t *testing.T) {
	todo := &Todo{
		ID:       3,
		Title:    "Toggle Test",
		Content:  "Testing status toggling",
		Finished: false,
	}

	todo.MarkFinished()
	assert.True(t, todo.Finished)

	todo.MarkUnfinished()
	assert.False(t, todo.Finished)

	todo.MarkFinished()
	assert.True(t, todo.Finished)
}
