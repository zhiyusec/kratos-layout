package repo

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v3/errors"
)

const (
	errorReasonTodoNotFound        = "TODO_NOT_FOUND"
	errorReasonTodoInvalidArgument = "TODO_INVALID_ARGUMENT"
)

var (
	ErrTodoNotFound        = errors.NotFound(errorReasonTodoNotFound, "todo not found")
	ErrTodoInvalidArgument = errors.BadRequest(errorReasonTodoInvalidArgument, "invalid todo argument")
)

// Todo is a domain model.
type Todo struct {
	ID         int64
	Title      string
	Content    string
	Completed  bool
	CreateTime time.Time
	UpdateTime time.Time
}

// TodoRepo provides data access for Todo.
type TodoRepo struct {
	mu     sync.RWMutex
	nextID int64
	todos  map[int64]*Todo
}

func NewTodoRepo() *TodoRepo {
	return &TodoRepo{
		nextID: 1,
		todos:  make(map[int64]*Todo),
	}
}

func (r *TodoRepo) FindByID(_ context.Context, id int64) (*Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	todo, ok := r.todos[id]
	if !ok {
		return nil, ErrTodoNotFound
	}
	return clone(todo), nil
}

func (r *TodoRepo) List(_ context.Context, offset, limit int) ([]*Todo, error) {
	if offset < 0 || limit <= 0 {
		return nil, ErrTodoInvalidArgument
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	todos := make([]*Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, clone(todo))
	}
	slices.SortFunc(todos, func(a, b *Todo) int {
		return cmp.Compare(a.ID, b.ID)
	})
	if offset >= len(todos) {
		return []*Todo{}, nil
	}
	end := offset + limit
	if end > len(todos) {
		end = len(todos)
	}
	return todos[offset:end], nil
}

func (r *TodoRepo) Create(_ context.Context, todo *Todo) (*Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	todo = clone(todo)
	todo.ID = r.nextID
	todo.CreateTime = now
	todo.UpdateTime = now
	r.todos[todo.ID] = clone(todo)
	r.nextID++
	return clone(todo), nil
}

func (r *TodoRepo) Update(_ context.Context, todo *Todo) (*Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	current, ok := r.todos[todo.ID]
	if !ok {
		return nil, ErrTodoNotFound
	}
	updated := clone(todo)
	updated.CreateTime = current.CreateTime
	updated.UpdateTime = time.Now()
	r.todos[updated.ID] = clone(updated)
	return clone(updated), nil
}

func (r *TodoRepo) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.todos[id]; !ok {
		return ErrTodoNotFound
	}
	delete(r.todos, id)
	return nil
}

func clone(todo *Todo) *Todo {
	if todo == nil {
		return nil
	}
	return &Todo{
		ID:         todo.ID,
		Title:      todo.Title,
		Content:    todo.Content,
		Completed:  todo.Completed,
		CreateTime: todo.CreateTime,
		UpdateTime: todo.UpdateTime,
	}
}
