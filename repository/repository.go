package repository

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user interface{}) error
	GetByID(ctx context.Context, id interface{}, result interface{}) error
	List(ctx context.Context, filter interface{}, results interface{}) error
	Update(ctx context.Context, id interface{}, update interface{}) error
	Delete(ctx context.Context, id interface{}) error
	Count(ctx context.Context) (int64, error)
}
