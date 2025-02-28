package mocks

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/slon/shad-go/coverme/models"
)

type MockStorage struct {
	mock.Mock
}

func (mockStorage *MockStorage) AddTodo(title, content string) (*models.Todo, error) {
	args := mockStorage.Called(title, content)
	todo, _ := args.Get(0).(*models.Todo)
	return todo, args.Error(1)
}

func (mockStorage *MockStorage) GetTodo(id models.ID) (*models.Todo, error) {
	args := mockStorage.Called(id)
	todo, _ := args.Get(0).(*models.Todo)
	return todo, args.Error(1)
}

func (mockStorage *MockStorage) GetAll() ([]*models.Todo, error) {
	args := mockStorage.Called()
	todos, _ := args.Get(0).([]*models.Todo)
	return todos, args.Error(1)
}

func (mockStorage *MockStorage) FinishTodo(id models.ID) error {
	args := mockStorage.Called(id)
	return args.Error(0)
}
