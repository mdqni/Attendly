package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	group1 "github.com/mdqni/Attendly/proto/gen/go/group/v1"
	userv1 "github.com/mdqni/Attendly/proto/gen/go/user/v1"
	errs "github.com/mdqni/Attendly/shared/errs"
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

	return &PostgresRepo{db: pool}, nil
}

func (r *PostgresRepo) CreateGroup(ctx context.Context, groupName string, department string, year int) (*group1.Group, error) {
	op := "groups.storage.postgres.CreateGroup"

	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO "group".groups (name, department, year)
		VALUES ($1, $2, $3)
		RETURNING id
	`, groupName, department, year).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errs.ErrGroupAlreadyExists
		}
		return nil, fmt.Errorf("%s: failed to create group: %w", op, err)
	}

	return &group1.Group{
		Id:         id,
		Name:       groupName,
		Department: department,
		Year:       int32(year),
	}, nil
}

func (r *PostgresRepo) AddUserToGroup(ctx context.Context, groupID string, userID string) (bool, error) {
	op := "groups.storage.postgres.AddUserToGroup"
	const query = `
		INSERT INTO "group".user_groups (user_id, group_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING;
	`

	_, err := r.db.Exec(ctx, query, userID, groupID)
	if err != nil {
		log.Println("op", op, "Error: ", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return false, errs.ErrUserOrGroupNotExists
		}
		return false, fmt.Errorf("failed to add user to group: %w", err)
	}
	return true, nil
}

func (r *PostgresRepo) RemoveUserFromGroup(ctx context.Context, groupID, userID string) (bool, error) {
	const query = `DELETE FROM "group".user_groups
		WHERE user_id = $1 AND group_id = $2;`
	row, err := r.db.Exec(ctx, query, userID, groupID)
	if err != nil {
		return false, fmt.Errorf("failed to remove user from group: %w", err)
	}
	if row.RowsAffected() == 0 {
		return false, errs.ErrNoUserInGroup
	}
	return true, nil
}

func (r *PostgresRepo) GetGroup(ctx context.Context, groupID string) (*group1.Group, error) {
	var group group1.Group
	err := r.db.QueryRow(ctx, `SELECT id, name, department, year FROM "group".groups WHERE id = $1`, groupID).
		Scan(&group.Id, &group.Name, &group.Department, &group.Year)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrGroupNotFound
		}
		return nil, err
	}

	return &group1.Group{
		Id:         group.Id,
		Name:       group.Name,
		Department: group.Department,
		Year:       group.Year,
	}, nil
}

func (r *PostgresRepo) ListUsersInGroup(ctx context.Context, groupID string) ([]*userv1.User, error) {
	const query = `
		SELECT u.id, u.name, u.barcode, r.name
		FROM "group".user_groups ug
		JOIN "auth".users u ON u.id = ug.user_id
		JOIN "auth".roles r ON r.id = u.role_id
		WHERE ug.group_id = $1;
	`

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("ListUsersInGroup query error: %w", err)
	}
	defer rows.Close()

	var users []*userv1.User

	for rows.Next() {
		var user userv1.User
		var role string
		if err := rows.Scan(&user.Id, &user.Name, &user.Barcode, &role); err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		user.Role = role
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(users) == 0 {
		return nil, errs.ErrGroupIsEmpty
	}

	return users, nil
}

func (r *PostgresRepo) IsInGroup(ctx context.Context, groupID string, userId string) (bool, error) {
	const query = `
		SELECT EXISTS (
		SELECT 1
		FROM "group".user_groups
		where user_id = $1 AND group_id = $2
		)`
	var exists bool
	err := r.db.QueryRow(ctx, query, userId, groupID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("IsInGroup query error: %w", err)
	}
	return exists, nil
}
