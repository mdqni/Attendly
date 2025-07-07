package repository

import (
	"context"
	"github.com/mdqni/Attendly/services/auth/internal/domain/model"
	"time"
)

type AuthRepository interface {
	SaveUser(ctx context.Context, user model.UserWithPassword) error
	GetUserByBarcode(ctx context.Context, barcode string) (*model.UserWithPassword, error)
	GetPermissions(ctx context.Context, userID string) ([]string, error)
	SaveRefreshToken(ctx context.Context, token string, userID string, expiresAt time.Time) error
	ValidateRefreshToken(ctx context.Context, token string) (string, error)
	DeleteRefreshToken(ctx context.Context, token string) error
}
