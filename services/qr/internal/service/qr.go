package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	redis2 "github.com/mdqni/Attendly/services/qr/internal/repository/redis"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type qrService struct {
	rdb *redis2.RedisCache
}

func NewQrService(client *redis2.RedisCache) QrService {
	return &qrService{rdb: client}
}

const qrPrefix = "qr:"
const qrTTL = 5 * time.Minute

type qrPayload struct {
	LessonID  string `json:"lesson_id"`
	TeacherID string `json:"teacher_id"`
	ExpiresAt int64  `json:"expires_unix"`
}

func (q qrService) GenerateQR(ctx context.Context, lessonID string, teacherID string, unix int64) (string, error) {
	token := uuid.NewString()
	key := qrPrefix + token

	payload := qrPayload{
		LessonID:  lessonID,
		TeacherID: teacherID,
		ExpiresAt: unix,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal qr payload: %w", err)
	}

	err = q.rdb.Set(ctx, key, data, qrTTL)
	if err != nil {
		return "", fmt.Errorf("failed to store QR in redis: %w", err)
	}

	return token, nil
}

func (q qrService) ValidateQR(ctx context.Context, token string) (string, error) {
	key := qrPrefix + token

	data, err := q.rdb.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", status.Error(codes.NotFound, "QR code expired or invalid")
		}
		return "", fmt.Errorf("failed to validate QR: %w", err)
	}

	var payload qrPayload
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return "", fmt.Errorf("failed to unmarshal QR data: %w", err)
	}

	return payload.LessonID, nil
}
