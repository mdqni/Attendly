package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mdqni/Attendly/shared/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	errPkg "github.com/mdqni/Attendly/shared/errs"
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

func (r *PostgresRepo) GetUserByID(ctx context.Context, id string) (*userv1.User, error) {
	query := `
	SELECT id, name, email, avatar_url
	FROM "user".user_profiles
	WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	var u userv1.User
	err := row.Scan(&u.Id, &u.Name, &u.Email, &u.AvatarUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errPkg.ErrUserNotFound
		}
		return nil, fmt.Errorf("repo.GetUserByID: %w", err)
	}
	return &u, nil
}

func (r *PostgresRepo) GetUsers(ctx context.Context, page, limit int32) ([]*userv1.User, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	query := `
	SELECT id, name, email, avatar_url
	FROM "user".user_profiles
	ORDER BY name
	OFFSET $1 LIMIT $2
	`
	rows, err := r.db.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("repo.GetUsers: %w", err)
	}
	defer rows.Close()

	var users []*userv1.User
	for rows.Next() {
		var u userv1.User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.AvatarUrl); err != nil {
			return nil, fmt.Errorf("repo.GetUsers scan: %w", err)
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo.GetUsers rows: %w", err)
	}
	return users, nil
}
func (r *PostgresRepo) GetAllUsers(ctx context.Context, page int, limit int) ([]*userv1.User, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	query := `
	SELECT id, name, email, avatar_url
	FROM "user".user_profiles
	ORDER BY name
	OFFSET $1 LIMIT $2
	`
	rows, err := r.db.Query(ctx, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("repo.GetUsers: %w", err)
	}
	defer rows.Close()

	var users []*userv1.User
	for rows.Next() {
		var u userv1.User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.AvatarUrl); err != nil {
			return nil, fmt.Errorf("repo.GetUsers scan: %w", err)
		}
		users = append(users, &u)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("repo.GetUsers rows: %w", rows.Err())
	}
	return users, nil
}

func (r *PostgresRepo) UpdateUser(ctx context.Context, u *userv1.User) (*userv1.User, error) {
	query := `
	UPDATE "user".user_profiles
	SET name = $2, email = $3, avatar_url = $4
	WHERE id = $1
	RETURNING id, name, email, avatar_url
	`
	row := r.db.QueryRow(ctx, query, u.Id, u.Name, u.Email, u.AvatarUrl)
	var updated userv1.User
	if err := row.Scan(&updated.Id, &updated.Name, &updated.Email, &updated.AvatarUrl); err != nil {
		if err == sql.ErrNoRows {
			return nil, errPkg.ErrUserNotFound
		}
		return nil, fmt.Errorf("repo.UpdateUser: %w", err)
	}
	return &updated, nil
}

func (r *PostgresRepo) DeleteUser(ctx context.Context, id string) error {
	cmd, err := r.db.Exec(ctx, `DELETE FROM "user".user_profiles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("repo.DeleteUser: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return errPkg.ErrUserNotFound
	}
	return nil
}

func (r *PostgresRepo) IsUserInGroup(ctx context.Context, userID, groupID string) (bool, error) {
	panic("implement me")
}

func (r *PostgresRepo) CreateUser(ctx context.Context, u *domain.User) (*userv1.User, error) {
	query := `
	INSERT INTO "user".user_profiles (id, email, name)
	VALUES ($1, $2, $3)
	RETURNING id, name, email, avatar_url
	`

	row := r.db.QueryRow(ctx, query, u.ID, u.Email, u.Name)

	var result userv1.User
	if err := row.Scan(&result.Id, &result.Name, &result.Email, &result.Name); err != nil {
		return nil, fmt.Errorf("repo.CreateUser: %w", err)
	}
	return &result, nil
}
