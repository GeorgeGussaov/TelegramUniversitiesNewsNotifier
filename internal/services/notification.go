package notification_service

import (
	"context"
	"encoding/json"
	"log"
	"notification-bot/internal/kafka"
	"notification-bot/internal/models"
	"notification-bot/internal/notifier"
	"notification-bot/internal/repository"
)

type NotificationService struct {
	consumer *kafka.Consumer
	notifier *notifier.TelegramNotifier
	newsRepo *repository.NewsRepository 
}

func NewNotificationService(consumer *kafka.Consumer, notifier *notifier.TelegramNotifier, newsRepos *repository.NewsRepository) *NotificationService {
	return &NotificationService{
		consumer: consumer,
		notifier: notifier,
		newsRepo: newsRepos,
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
		check, err := s.newsRepo.IsNewsSent(news.Link)	//при перезапуске кода надо проверять в бд уже отправленные новости
		if err != nil{
			log.Print(err)
			return
		}
		if check{
			log.Printf("news already sent: %s", news.Link)
			return
		}
		err = s.notifier.Notify(news)
		if err != nil {
			log.Printf("failed to send tg message: %v", err)
			return
		}
		s.newsRepo.Add(&news)
	})
}