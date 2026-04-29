package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func (c DBConfig) GetHost() string {
	return c.Host
}

func (c DBConfig) GetPort() int {
	return c.Port
}

func (c DBConfig) GetUser() string {
	return c.User
}

func (c DBConfig) GetPassword() string {
	return c.Password
}

func (c DBConfig) GetName() string {
	return c.Name
}
type KafkaConfig struct{
	KafkaBrokers string
	KafkaTopic string
	KafkaGroupID string
}
type Config struct {
	DB            DBConfig
	TelegramToken string
	TelegramChatID int64
	KafkaCfg KafkaConfig
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load("internal/config/config.env"); err != nil {
		log.Println(err)
	}
	chatID := os.Getenv("CHANNEL_CHAT_ID")
	intChatID, _ := strconv.ParseInt(chatID, 10, 64)
	return &Config{
		DB: DBConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "UniversitiesNews",
		},
		KafkaCfg: KafkaConfig{
			KafkaBrokers: os.Getenv("KAFKA_BROKER"),
			KafkaTopic: os.Getenv("KAFKA_TOPIC"),
			KafkaGroupID: os.Getenv("KAFKA_GROUP_ID"),
		},
		TelegramChatID: intChatID,
		TelegramToken: os.Getenv("API_TOKEN"), //passengers_bot
	}, nil
}
