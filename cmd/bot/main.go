package main

import (
	"notification-bot/internal/config"
	"notification-bot/internal/kafka"
	"notification-bot/internal/notifier"
	"notification-bot/internal/repository"
	notification_service "notification-bot/internal/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := repository.NewDB(cfg.DB)
	if err != nil {
	 	log.Fatalf("Failed to connect to database: %v", err)
	}
	newsRepo, err := repository.NewNewsRepository(db)
	if err != nil{
		log.Fatalf("Failed to create news repository: %v", err)
	}
	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Printf("failed create bot: %v", err)
	}


	tgNotifier := notifier.NewTelegramNotifier(botAPI, cfg.TelegramChatID)

	log.Println(cfg.KafkaCfg.KafkaBrokers)

	consumer := kafka.NewConsumer(
		[]string{cfg.KafkaCfg.KafkaBrokers},
		cfg.KafkaCfg.KafkaTopic,
		cfg.KafkaCfg.KafkaGroupID,
	)

	notificationService := notification_service.NewNotificationService(consumer, tgNotifier, newsRepo)

	go notificationService.Start()

	go StartBot(botAPI)

	log.Println("service started")

	select {}
}

func StartBot(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {
			log.Println("chat_id:", update.Message.Chat.ID)
			continue
		}

		if update.ChannelPost != nil {
			log.Println("channel chat_id:", update.ChannelPost.Chat.ID)
			continue
		}
	}
}