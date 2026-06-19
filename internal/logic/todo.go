package logic

import (
	"context"
	"strings"

	"github.com/zhiyusec/kratos-layout/internal/repo"
)

type TodoLogic struct {
	repo *repo.TodoRepo
}

func NewTodoLogic(repo *repo.TodoRepo) *TodoLogic {
	return &TodoLogic{repo: repo}
}

func (l *TodoLogic) Create(ctx context.Context, todo *repo.Todo) (*repo.Todo, error) {
	if err := validate(todo); err != nil {
		return nil, err
	}
	return l.repo.Create(ctx, todo)
}

func (l *TodoLogic) Get(ctx context.Context, id int64) (*repo.Todo, error) {
	return l.repo.FindByID(ctx, id)
}

func (l *TodoLogic) List(ctx context.Context, offset, limit int) ([]*repo.Todo, error) {
	return l.repo.List(ctx, offset, limit)
}

func (l *TodoLogic) Update(ctx context.Context, todo *repo.Todo) (*repo.Todo, error) {
	if todo == nil || todo.ID <= 0 {
		return nil, repo.ErrTodoInvalidArgument
	}
	if err := validate(todo); err != nil {
		return nil, err
	}
	return l.repo.Update(ctx, todo)
}

func (l *TodoLogic) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return repo.ErrTodoInvalidArgument
	}
	return l.repo.Delete(ctx, id)
}

func validate(todo *repo.Todo) error {
	if todo == nil || strings.TrimSpace(todo.Title) == "" {
		return repo.ErrTodoInvalidArgument
	}
	return nil
}
