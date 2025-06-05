package config

import (
	"os"
)

// Config holds bot configuration
type Config struct {
	BotToken     string
	MongoURI     string
	MongoDBName  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return Config{}, fmt.Errorf("BOT_TOKEN not set")
	}
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return Config{}, fmt.Errorf("MONGO_URI not set")
	}
	mongoDBName := os.Getenv("MONGO_DB_NAME")
	if mongoDBName == "" {
		mongoDBName = "telegram_bot"
	}
	return Config{
		BotToken:    botToken,
		MongoURI:    mongoURI,
		MongoDBName: mongoDBName,
	}, nil
}
