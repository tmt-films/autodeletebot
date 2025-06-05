package main

import (
	"log"
	"os"

	"github.com/tmt-films/autodeletebot/config"
	"github.com/tmt-films/autodeletebot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}

	// Set up update configuration
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Initialize database and handlers
	handlers.InitDB(cfg)
	handlers.InitScheduler(bot)

	// Start receiving updates
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		handlers.HandleUpdate(bot, update)
	}
}
