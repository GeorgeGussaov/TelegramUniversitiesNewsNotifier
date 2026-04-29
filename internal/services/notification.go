package notification_service

import (
	"context"
	"encoding/json"
	"log"
	"notification-bot/internal/kafka"
	"notification-bot/internal/models"
	"notification-bot/internal/notifier"
)

type NotificationService struct {
	consumer *kafka.Consumer
	notifier *notifier.TelegramNotifier 
}

func NewNotificationService(consumer *kafka.Consumer, notifier *notifier.TelegramNotifier) *NotificationService {
	return &NotificationService{
		consumer: consumer,
		notifier: notifier,
	}
}

func (s *NotificationService) Start() {
	ctx := context.Background()

	
	s.consumer.ReadMessages(ctx, func(msg []byte) {

		var news models.News

		err := json.Unmarshal(msg, &news)
		if err != nil {
			log.Printf("failed to parse json: %v", err)
			return
		}

		err = s.notifier.Notify(news)
		if err != nil {
			log.Printf("failed to send tg message: %v", err)
		}
	})
}