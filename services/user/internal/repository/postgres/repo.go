package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	"github.com/mdqni/Attendly/services/user/internal/utils/passwordUtils"
	"log"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func New(connString string) (*PostgresRepo, error) {
	const op = "storage.postgres.New"

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("%s: failed to connect to database: %w", op, err)
	}

	_, err = pool.Exec(ctx, `SET search_path TO "user"`)
	if err != nil {
		log.Println(err)
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

		`CREATE TABLE IF NOT EXISTS roles (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			barcode TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			role_id INT REFERENCES roles(id)
		);
		
		CREATE TABLE IF NOT EXISTS permissions (
			id SERIAL PRIMARY KEY,
			action TEXT UNIQUE NOT NULL
		);
		
		
		CREATE TABLE IF NOT EXISTS role_permissions (
			role_id INT REFERENCES roles(id) ON DELETE CASCADE,
			permission_id INT REFERENCES permissions(id) ON DELETE CASCADE,
			PRIMARY KEY (role_id, permission_id)
		);
		
		INSERT INTO roles (name) VALUES ('admin'), ('teacher'), ('student')
		ON CONFLICT DO NOTHING;`,
	}

	for _, stmt := range statements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return nil, fmt.Errorf("%s: failed to execute statement: %w\nSQL: %s", op, err, stmt)
		}
	}

	return &PostgresRepo{db: pool}, nil
}
func (r *PostgresRepo) SaveUser(ctx context.Context, user *userv1.User) error {
	op := "storage.postgres.saveUser"
	var roleID int
	err := r.db.QueryRow(ctx, `SELECT id FROM roles WHERE name = $1`, user.Role).Scan(&roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}
	hashed, err := passwordUtils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashed

	_, err = r.db.Exec(ctx, `
		INSERT INTO users (id, name, barcode, password, role_id)
		VALUES ($1, $2, $3, $4, $5)
	`, user.Id, user.Name, user.Barcode, user.Password, roleID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return fmt.Errorf("barcode already exists: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return err
}

func (r *PostgresRepo) GetUserById(ctx context.Context, id string) (*userv1.User, error) {
	const query = `
		SELECT id, name, barcode, role_id, password
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var user userv1.User
	err := row.Scan(&user.Id, &user.Name, &user.Barcode, &user.Role, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var role string
	err = r.db.QueryRow(ctx, `SELECT name FROM roles WHERE id = $1`, user.Role).Scan(&role)
	if err != nil {
		return nil, err
	}
	user.Role = role
	return &user, nil
}

func (r *PostgresRepo) GetUserByBarcode(ctx context.Context, barcode string) (*userv1.User, error) {
	const query = `
		SELECT id, name, barcode, role_id, password
		FROM users
		WHERE barcode = $1
	`

	row := r.db.QueryRow(ctx, query, barcode)

	var user userv1.User
	err := row.Scan(&user.Id, &user.Name, &user.Barcode, &user.Role, &user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var role string
	err = r.db.QueryRow(ctx, `SELECT name FROM roles WHERE id = $1`, user.Role).Scan(&role)
	if err != nil {
		return nil, err
	}
	user.Role = role
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

func (r *PostgresRepo) GetPermissions(ctx context.Context, userID string) ([]string, error) {
	query := `
	SELECT p.action
	FROM users u
	JOIN roles r ON u.role_id = r.id
	JOIN role_permissions rp ON rp.role_id = r.id
	JOIN permissions p ON p.id = rp.permission_id
	WHERE u.id = $1;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}
	log.Println("Permissions:", perms)
	return perms, rows.Err()
}
