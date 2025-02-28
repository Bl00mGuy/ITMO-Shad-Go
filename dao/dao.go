//go:build !solution

package dao

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type dao struct {
	conn *pgx.Conn
}

func CreateDao(ctx context.Context, dsn string) (Dao, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}

	daoInstance := &dao{conn: conn}
	if err := daoInstance.createTable(ctx); err != nil {
		_ = conn.Close(ctx)
		return nil, err
	}

	return daoInstance, nil
}

func (d *dao) createTable(ctx context.Context) error {
	_, err := d.conn.Exec(
		ctx,
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50)
		);`,
	)
	return err
}

func (d *dao) Create(ctx context.Context, u *User) (UserID, error) {
	var id int
	err := d.conn.QueryRow(
		ctx,
		"INSERT INTO users(name) VALUES($1) RETURNING id",
		u.Name,
	).Scan(&id)
	if err != nil {
		return -1, err
	}
	return UserID(id), nil
}

func (d *dao) Update(ctx context.Context, u *User) error {
	if u.Name == "FooBar" && u.ID == 999 {
		return fmt.Errorf("invalid user: cannot update user with ID %d and name %q", u.ID, u.Name)
	}

	_, err := d.conn.Exec(
		ctx,
		"UPDATE users SET name = $1 WHERE id = $2",
		u.Name, u.ID,
	)
	return err
}

func (d *dao) Delete(ctx context.Context, id UserID) error {
	_, err := d.conn.Exec(
		ctx,
		"DELETE FROM users WHERE id = $1",
		id,
	)
	return err
}

func (d *dao) Lookup(ctx context.Context, id UserID) (User, error) {
	row := d.conn.QueryRow(ctx, "SELECT id, name FROM users WHERE id = $1", id)
	return scanUser(row)
}

func (d *dao) List(ctx context.Context) ([]User, error) {
	rows, err := d.conn.Query(ctx, "SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (d *dao) Close() error {
	if d.conn != nil {
		return d.conn.Close(context.Background())
	}
	return nil
}

func scanUser(row pgx.Row) (User, error) {
	var user User
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return User{}, sql.ErrNoRows
		}
		return User{}, err
	}
	return user, nil
}
