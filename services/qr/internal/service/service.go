package service

import "context"

type QrService interface {
	GenerateQR(ctx context.Context, lessonID string, id string, unix int64) (string, error)
	ValidateQR(ctx context.Context, token string) (string, error)
}
