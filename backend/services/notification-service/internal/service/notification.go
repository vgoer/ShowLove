package service

import (
	"context"
	"fmt"
	"log"

	"showlove/pkg/events"
	"showlove/services/notification-service/internal/model"
	"showlove/services/notification-service/internal/repository"
)

type NotificationService struct {
	repo repository.DeviceRepository
	sub  events.Subscriber
}

func NewNotificationService(repo repository.DeviceRepository, sub events.Subscriber) *NotificationService {
	return &NotificationService{repo: repo, sub: sub}
}

func (s *NotificationService) RegisterDevice(ctx context.Context, userID, token, platform string) (*model.DeviceToken, error) {
	dt := &model.DeviceToken{
		UserID:   userID,
		Token:    token,
		Platform: platform,
	}
	if err := s.repo.Upsert(ctx, dt); err != nil {
		return nil, fmt.Errorf("notification: register device: %w", err)
	}
	return dt, nil
}

// SendPush sends a push notification (placeholder for FCM/APNs integration).
func (s *NotificationService) SendPush(ctx context.Context, userID, title, body string) error {
	tokens, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	for _, dt := range tokens {
		log.Printf("[notification] Push to user=%s device=%s platform=%s: %s - %s",
			userID, dt.Token[:8]+"...", dt.Platform, title, body)
		// TODO: Call FCM/APNs via push adapter
	}
	return nil
}
