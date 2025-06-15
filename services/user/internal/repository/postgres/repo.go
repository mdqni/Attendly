package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func New(connString string) (*PostgresRepo, error) {
	const op = "storage.postgres.New"

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to database: %w", op, err)
	}

	_, err = pool.Exec(ctx, `SET search_path TO "user"`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to set search_path: %w", op, err)
	}

	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`,

		`SET search_path TO "user";`,

		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_status') THEN
				CREATE TYPE role_status AS ENUM ('student', 'teacher', 'admin');
			END IF;
		END
		$$;`,

		`CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    barcode TEXT UNIQUE NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('student', 'teacher', 'admin')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);`,
	}

	for _, stmt := range statements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return nil, fmt.Errorf("%s: failed to execute statement: %w\nSQL: %s", op, err, stmt)
		}
	}

	return &PostgresRepo{db: pool}, nil
}
func (r *PostgresRepo) SaveUser(ctx context.Context, user *userv1.User) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (id, name, barcode, role)
		VALUES ($1, $2, $3, $4)
	`, user.Id, user.Name, user.Barcode, user.Role)
	return err
}

func (r *PostgresRepo) GetUserByID(ctx context.Context, id string) (*userv1.User, error) {
	const query = `
		SELECT id, name, barcode, role
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var user userv1.User
	err := row.Scan(&user.Id, &user.Name, &user.Barcode, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepo) CheckUserInGroup(ctx context.Context, userID, groupID string) (bool, error) {
	return true, nil // заглушка
}
