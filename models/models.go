package models

// ChatSettings holds chat-specific settings
type ChatSettings struct {
	ChatID      int64  `bson:"chat_id"`
	Timer       int    `bson:"timer"`
	NoMedia     int    `bson:"no_media"`
	BannedWords string `bson:"banned_words"`
}
