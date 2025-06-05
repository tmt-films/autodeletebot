package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/YourUsername/telegram-auto-delete-bot/config"
	"github.com/YourUsername/telegram-auto-delete-bot/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var mongoClient *mongo.Client
var settingsCollection *mongo.Collection
var messagesCollection *mongo.Collection

// InitDB initializes the MongoDB connection and collections
func InitDB(cfg config.Config) {
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	var err error
	mongoClient, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Initialize collections
	db := mongoClient.Database(cfg.MongoDBName)
	settingsCollection = db.Collection("settings")
	messagesCollection = db.Collection("messages")

	// Create TTL index for automatic message deletion
	ttlIndex := mongo.IndexModel{
		Keys: bson.M{"timestamp": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, err = messagesCollection.Indexes().CreateOne(context.Background(), ttlIndex)
	if err != nil {
		log.Fatalf("Failed to create TTL index: %v", err)
	}
}

// InitScheduler initializes a simple scheduler for non-TTL deletions
func InitScheduler(bot *tgbotapi.BotAPI) {
	go func() {
		for {
			time.Sleep(time.Minute)
			// MongoDB TTL handles time-based deletion, but this can be used for other cleanup tasks
		}
	}()
}

// HandleUpdate processes incoming Telegram updates
func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	isAdmin := isAdmin(bot, update.Message.From.ID, chatID)

	// Handle commands
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			handleStart(bot, update)
		case "help":
			handleHelp(bot, update)
		case "settimer":
			if isAdmin {
				handleSetTimer(bot, update)
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Only admins can use /settimer."))
			}
		case "banword":
			if isAdmin {
				handleBanWord(bot, update)
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Only admins can use /banword."))
			}
		case "nomedia":
			if isAdmin {
				handleNoMedia(bot, update)
			} else {
				bot.Send(tgbotapi.NewMessage(chatID, "Only admins can use /nomedia."))
			}
		}
		return
	}

	// Check message for deletion
	checkMessage(bot, update.Message)
}

// isAdmin checks if the user is an admin in the chat
func isAdmin(bot *tgbotapi.BotAPI, userID int64, chatID int64) bool {
	member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chatID,
			UserID: userID,
		},
	})
	if err != nil {
		log.Printf("Failed to get chat member: %v", err)
		return false
	}
	return member.IsAdministrator() || member.IsCreator()
}

// handleStart initializes the chat and sends a welcome message
func handleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	// Initialize chat settings in MongoDB if not present
	_, err := settingsCollection.UpdateOne(
		context.Background(),
		bson.M{"chat_id": chatID},
		bson.M{"$setOnInsert": bson.M{
			"timer":       0,
			"no_media":    0,
			"banned_words": "",
		}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to initialize chat settings: %v", err)
		bot.Send(tgbotapi.NewMessage(chatID, "Error initializing bot. Please try again."))
		return
	}

	// Send welcome message
	msg := tgbotapi.NewMessage(chatID, "Welcome to the Auto-Delete Bot! I can delete messages based on time, content, or media type. Use /help to see available commands.")
	bot.Send(msg)
}

// handleHelp sends a list of available commands
func handleHelp(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	helpText := `Available commands:
- /start - Initialize the bot and see this welcome message.
- /help - Show this help message.
- /settimer <minutes>m - Set a timer to delete messages after a specified time (e.g., /settimer 5m) [Admin only].
- /banword <word> - Delete messages containing the specified word (e.g., /banword spam) [Admin only].
- /nomedia - Delete all photos, videos, and stickers [Admin only].`
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
	bot.Send(msg)
}

// handleSetTimer sets the auto-delete timer for the chat
func handleSetTimer(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	args := update.Message.CommandArguments()
	var minutes int
	_, err := fmt.Sscanf(args, "%dm", &minutes)
	if err != nil || minutes <= 0 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /settimer <minutes>m (e.g., /settimer 5m)"))
		return
	}

	_, err = settingsCollection.UpdateOne(
		context.Background(),
		bson.M{"chat_id": update.Message.Chat.ID},
		bson.M{"$set": bson.M{"timer": minutes * 60}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to set timer: %v", err)
		return
	}
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Messages will be deleted after %d minutes.", minutes)))
}

// handleBanWord adds a banned word for the chat
func handleBanWord(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	word := update.Message.CommandArguments()
	if word == "" {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Usage: /banword <word>"))
		return
	}

	var settings models.ChatSettings
	err := settingsCollection.FindOne(context.Background(), bson.M{"chat_id": update.Message.Chat.ID}).Decode(&settings)
	if err == mongo.ErrNoDocuments {
		settings = models.ChatSettings{ChatID: update.Message.Chat.ID, BannedWords: word}
	} else if err != nil {
		log.Printf("Failed to get banned words: %v", err)
		return
	} else {
		if settings.BannedWords == "" {
			settings.BannedWords = word
		} else {
			settings.BannedWords = strings.Join(append(strings.Split(settings.BannedWords, ","), word), ",")
		}
	}

	_, err = settingsCollection.UpdateOne(
		context.Background(),
		bson.M{"chat_id": update.Message.Chat.ID},
		bson.M{"$set": bson.M{"banned_words": settings.BannedWords}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to set banned word: %v", err)
		return
	}
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Banned word added: %s", word)))
}

// handleNoMedia toggles media deletion
func handleNoMedia(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	_, err := settingsCollection.UpdateOne(
		context.Background(),
		bson.M{"chat_id": update.Message.Chat.ID},
		bson.M{"$set": bson.M{"no_media": 1}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to set no_media: %v", err)
		return
	}
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Media messages will be deleted."))
}

// checkMessage checks if a message should be deleted
func checkMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var settings models.ChatSettings
	err := settingsCollection.FindOne(context.Background(), bson.M{"chat_id": message.Chat.ID}).Decode(&settings)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Failed to get settings: %v", err)
		return
	}

	// Check for banned words
	if settings.BannedWords != "" {
		words := strings.Split(settings.BannedWords, ",")
		for _, word := range words {
			if word != "" && strings.Contains(strings.ToLower(message.Text), strings.ToLower(word)) {
				bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
				return
			}
		}
	}

	// Check for media
	if settings.NoMedia == 1 && (message.Photo != nil || message.Video != nil || message.Sticker != nil) {
		bot.Request(tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID))
		return
	}

	// Schedule for time-based deletion
	if settings.Timer > 0 {
		_, err := messagesCollection.InsertOne(context.Background(), bson.M{
			"chat_id":   message.Chat.ID,
			"message_id": message.MessageID,
			"timestamp": time.Now(),
			"expireAt":  time.Now().Add(time.Duration(settings.Timer) * time.Second),
		})
		if err != nil {
			log.Printf("Failed to store message: %v", err)
		}
	}
}
