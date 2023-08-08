package models

import (
	"context"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// GetAllUsers returns a slice of all users
func (m *DBModel) GetAllUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []*User

	query := `
		select
			id, last_name, first_name, email, created_at, updated_at
		from
			users
		order by
			last_name, first_name
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err = rows.Scan(
			&u.ID,
			&u.LastName,
			&u.FirstName,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

// GetOneUser returns one user by id
func (m *DBModel) GetOneUser(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u User

	query := `
		select
			id, last_name, first_name, email, created_at, updated_at
		from
			users
		where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&u.ID,
		&u.LastName,
		&u.FirstName,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// EditUser edits an existing user
func (m *DBModel) EditUser(u User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	stmt := `
		update users set
			first_name = $1,
			last_name = $2,
			email = $3,
			updated_at = $4
		where
			id = $5
			RETURNING id, first_name, last_name, email, created_at, updated_at`

	err := m.DB.QueryRowContext(ctx, stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		time.Now(),
		u.ID,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt)

	if err != nil {
		return User{}, err
	}
	return user, nil
}

// AddUser inserts a user into the database
func (m *DBModel) AddUser(u User, hash string) (User, error) {
	// func (m *DBModel) AddUser(u User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into users (first_name, last_name, email, password, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6)
		returning id`

	var id int

	err := m.DB.QueryRowContext(ctx, stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		hash,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return User{}, err
	}
	u.ID = id
	return u, nil
}

// DeleteUser deletes a user by id
func (m *DBModel) DeleteUser(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user, err := m.GetOneUser(id)

	if err != nil {
		return User{}, err
	}

	stmt := `delete from users where id = $1`

	_, err = m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return User{}, err
	}

	// stmt = "delete from tokens where user_id = ?"
	// _, err = m.DB.ExecContext(ctx, stmt, id)
	// if err != nil {
	// 	return err
	// }

	return user, nil
}
