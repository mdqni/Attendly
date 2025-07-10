package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mdqni/Attendly/services/auth/internal/domain/model"
	"github.com/mdqni/Attendly/shared/errs"
	"time"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

func New(connString string) (*PostgresRepo, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PostgresRepo{db: pool}, nil
}

func (r *PostgresRepo) SaveUser(ctx context.Context, user model.UserWithPassword) error {
	const op = "repo.SaveUser"
	var roleID int
	err := r.db.QueryRow(ctx, `SELECT id FROM roles WHERE name = $1`, user.Role).Scan(&roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO users (id, name, barcode, password, role_id)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.Name, user.Barcode, user.Password, roleID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errs.ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to save user: %w", err)
	}

	return err
}

func (r *PostgresRepo) GetUserByBarcode(ctx context.Context, barcode string) (*model.UserWithPassword, error) {
	row := r.db.QueryRow(ctx, `
        SELECT u.id,u.name,u.barcode,u.password,r.name
        FROM "auth".users u
        JOIN "auth".roles r ON u.role_id=r.id
        WHERE u.barcode=$1
    `, barcode)

	var u model.UserWithPassword
	if err := row.Scan(&u.ID, &u.Name, &u.Barcode, &u.Password, &u.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, fmt.Errorf("repo.GetUserByBarcode: %w", err)
	}
	return &u, nil
}

func (r *PostgresRepo) GetPermissions(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
        SELECT p.action
        FROM "auth".role_permissions rp
        JOIN "auth".permissions p ON rp.permission_id=p.id
        WHERE rp.role_id=(
            SELECT role_id FROM "auth".users WHERE id=$1
        )
    `, userID)
	if err != nil {
		return nil, fmt.Errorf("repo.GetPermissions: %w", err)
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var action string
		if err := rows.Scan(&action); err != nil {
			return nil, fmt.Errorf("repo.GetPermissions scan: %w", err)
		}
		perms = append(perms, action)
	}
	return perms, nil
}

func (r *PostgresRepo) SaveRefreshToken(ctx context.Context, token string, userID string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO "auth".refresh_tokens (token, user_id, expires_at)
        VALUES ($1,$2,$3)
    `, token, userID, expiresAt)
	if err != nil {
		return fmt.Errorf("repo.SaveRefreshToken: %w", err)
	}
	return nil
}

func (r *PostgresRepo) ValidateRefreshToken(ctx context.Context, token string) (string, error) {
	row := r.db.QueryRow(ctx, `
        SELECT user_id, expires_at FROM "auth".refresh_tokens WHERE token=$1
    `, token)

	var userID string
	var exp time.Time
	if err := row.Scan(&userID, &exp); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrTokenNotFound
		}
		return "", fmt.Errorf("repo.ValidateRefreshToken: %w", err)
	}
	if time.Now().After(exp) {
		return "", errs.ErrTokenExpired
	}
	return userID, nil
}

func (r *PostgresRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM "auth".refresh_tokens WHERE token=$1`, token)
	if err != nil {
		return fmt.Errorf("repo.DeleteRefreshToken: %w", err)
	}
	return nil
}
