package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	errs "github.com/mdqni/Attendly/shared/errs"
)

type PostgresRepo struct {
	db *pgxpool.Pool
}

type InternalUser struct {
	ID       string
	Name     string
	Barcode  string
	Password string
	Role     string
}

func New(connString string) (*PostgresRepo, error) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PostgresRepo{db: pool}, nil
}

func (r *PostgresRepo) GetUserById(ctx context.Context, id string) (*InternalUser, error) {
	query := `
		SELECT u.id, u.name, u.barcode, u.password, r.name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`

	row := r.db.QueryRow(ctx, query, id)
	var user InternalUser
	err := row.Scan(&user.ID, &user.Name, &user.Barcode, &user.Password, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetUserById error: %w", err)
	}

	return &user, nil
}

func (r *PostgresRepo) GetUserByBarcode(ctx context.Context, barcode string) (*InternalUser, error) {
	query := `
		SELECT u.id, u.name, u.barcode, u.password, r.name
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.barcode = $1
	`

	row := r.db.QueryRow(ctx, query, barcode)
	var user InternalUser
	err := row.Scan(&user.ID, &user.Name, &user.Barcode, &user.Password, &user.Role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetUserByBarcode error: %w", err)
	}

	return &user, nil
}

func (r *PostgresRepo) GetAllUsers(ctx context.Context) ([]*InternalUser, error) {
	query := `
		SELECT u.id, u.name, u.barcode, u.password, r.name
		FROM users u
		JOIN roles r ON u.role_id = r.id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers query error: %w", err)
	}
	defer rows.Close()

	var users []*InternalUser
	for rows.Next() {
		var user InternalUser
		if err := rows.Scan(&user.ID, &user.Name, &user.Barcode, &user.Password, &user.Role); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *PostgresRepo) IsUserInGroup(ctx context.Context, userID, groupID string) (bool, error) {
	return true, nil
}
