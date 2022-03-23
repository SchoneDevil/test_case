package db

import (
	"app/internal/user"
	pbUser "app/pb"
	"app/pkg/client/psql"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
)

type repository struct {
	client psql.Client
}

func (r *repository) Create(ctx context.Context, user *user.User) (string, error) {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	if err := r.client.QueryRow(ctx, query, user.Email, user.Password).Scan(&user.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return "", newErr
		}
		return "", err
	}

	return user.ID, nil
}

func (r *repository) FindAll(ctx context.Context) (u []*pbUser.User, err error) {
	query := `SELECT id, email, password FROM users;`

	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	users := make([]*pbUser.User, 0)

	for rows.Next() {
		var u user.User

		err = rows.Scan(&u.ID, &u.Email, &u.Password)
		if err != nil {
			return nil, err
		}

		users = append(users, &pbUser.User{
			Id:       u.ID,
			Email:    u.Email,
			Password: u.Password,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1 RETURNING id;`
	if err := r.client.QueryRow(ctx, query, id).Scan(id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return newErr
		}
		return err
	}

	return nil
}

func NewRepository(client psql.Client) user.Repository {
	return &repository{
		client: client,
	}
}
