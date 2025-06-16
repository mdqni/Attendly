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
);
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    action TEXT UNIQUE NOT NULL  -- e.g., 'user:create', 'user:read'
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INT REFERENCES roles(id),
    permission_id INT REFERENCES permissions(id),
    PRIMARY KEY(role_id, permission_id)
);

INSERT INTO roles (name) VALUES ('admin'), ('teacher'), ('student')
ON CONFLICT DO NOTHING;

INSERT INTO permissions (action) VALUES ('user:create'), ('user:read'), ('user:update')
ON CONFLICT DO NOTHING;

-- admin — все
INSERT INTO role_permissions (role_id, permission_id) VALUES 
(1,1), (1,2), (1,3)
ON CONFLICT DO NOTHING;

-- teacher
INSERT INTO role_permissions (role_id, permission_id) VALUES 
(2,2), (2,3)
ON CONFLICT DO NOTHING;

-- student
INSERT INTO role_permissions (role_id, permission_id) VALUES 
(3,2)
ON CONFLICT DO NOTHING;

`,
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
	return true, nil
}

func (r *PostgresRepo) HasPermission(ctx context.Context, userID string, action string) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM users u
			JOIN roles r ON r.name = u.role
			JOIN role_permissions rp ON rp.role_id = r.id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE u.id = $1 AND p.action = $2
		);
	`

	var hasPermission bool
	err := r.db.QueryRow(ctx, query, userID, action).Scan(&hasPermission)
	if err != nil {
		return false, fmt.Errorf("CheckPermission query error: %w", err)
	}

	return hasPermission, nil
}
