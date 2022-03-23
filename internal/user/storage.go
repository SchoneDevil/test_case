package user

import (
	user "app/pb"
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *User) (id string, err error)
	FindAll(ctx context.Context) (u []*user.User, err error)
	Delete(ctx context.Context, id string) error
}
