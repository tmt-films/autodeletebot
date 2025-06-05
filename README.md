Telegram Auto-Delete Bot
A Telegram bot written in Go that automatically deletes messages in group chats based on time, content, user, or media type. It uses MongoDB for storing chat settings and message data, with commands like /start and /help for user interaction. Ideal for keeping group chats clean and enforcing community rules.
Features

Time-Based Deletion: Deletes messages after a set time (e.g., /settimer 5m for 5 minutes).
Content-Based Deletion: Deletes messages containing banned words (e.g., /banword spam).
User-Specific Deletion: Deletes messages from non-admins or specific users.
Media-Based Deletion: Deletes photos, videos, or stickers (e.g., /nomedia).
Admin Controls: Restricts configuration commands (/settimer, /banword, /nomedia) to group admins.
Start Command: Initializes the bot and chat settings with /start.
Help Command: Lists all available commands with /help.
MongoDB Storage: Stores settings and messages with a TTL index for automatic deletion.

Prerequisites

Go 1.21 or higher
A Telegram bot token from BotFather
MongoDB instance (e.g., MongoDB Atlas, local MongoDB, or cloud provider)
MongoDB URI and database name

Setup

Clone the Repository:
git clone https://github.com/YourUsername/telegram-auto-delete-bot.git
cd telegram-auto-delete-bot


Install Dependencies:
go mod tidy


Set Environment Variables:Create a .env file or set the following:
export BOT_TOKEN="your-bot-token-here"
export MONGO_URI="mongodb+srv://<username>:<password>@<cluster>.mongodb.net/?retryWrites=true&w=majority"
export MONGO_DB_NAME="telegram_bot"


Run the Bot Locally:
go run main.go


Add the Bot to a Group:

Add the bot to your Telegram group and make it an admin with "Delete Messages" permission.
Use /setcommands in BotFather to register commands:start - Initialize the bot
help - Show available commands
settimer - Set auto-delete timer (e.g., /settimer 5m)
banword - Ban a word (e.g., /banword spam)
nomedia - Delete all media messages





Usage

/start: Initializes the bot and creates default chat settings. Responds with a welcome message.
/help: Lists all commands with descriptions.
/settimer <minutes>m: Sets a timer to delete messages after the specified time (e.g., /settimer 5m) [Admin only].
/banword <word>: Deletes messages containing the specified word (e.g., /banword spam) [Admin only].
/nomedia: Deletes all photos, videos, and stickers [Admin only].

MongoDB Setup

Create a MongoDB database (e.g., via MongoDB Atlas free tier or local MongoDB with docker run -d -p 27017:27017 mongo).
Get the connection URI (e.g., mongodb+srv://<username>:<password>@<cluster>.mongodb.net/?retryWrites=true&w=majority).
Set MONGO_URI and MONGO_DB_NAME in your environment.
The bot uses two collections:
settings: Stores chat-specific settings (chat_id, timer, no_media, banned_words).
messages: Stores messages with a TTL index for auto-deletion based on the expireAt field.



Deployment
To run the bot 24/7, deploy on a cloud platform like Heroku, AWS, or Render:

Set up a Go environment.
Configure BOT_TOKEN, MONGO_URI, and MONGO_DB_NAME as environment variables.
Ensure MongoDB is accessible (e.g., whitelist your server’s IP in MongoDB Atlas).
Use a persistent MongoDB instance (e.g., MongoDB Atlas for cloud hosting).
Example for Heroku:heroku create
heroku config:set BOT_TOKEN="your-bot-token" MONGO_URI="your-mongo-uri" MONGO_DB_NAME="telegram_bot"
git push heroku main



Testing

Add the bot to a private Telegram group and make it an admin.
Test commands:
/start: Should respond with a welcome message and create a settings document in MongoDB.
/help: Should list all commands.
/settimer 1m: Set a 1-minute deletion timer (admin only).
/banword test: Add "test" as a banned word (admin only).
/nomedia: Enable media deletion (admin only).
Send messages with banned words, media, or wait for the timer to verify deletions.


Check MongoDB collections (settings and messages) to confirm data storage.

Notes

Replace YourUsername: Update go.mod and README.md with your GitHub username.
Error Handling: The bot logs errors to the console. For production, consider adding a logging service (e.g., to a file or external service).
Security: Admin-only commands are restricted using isAdmin. Ensure the bot has "Delete Messages" permission in groups.
Scalability: MongoDB’s TTL index handles time-based deletions efficiently. For high-traffic groups, consider sharding or replica sets in MongoDB.
Enhancements: Add inline buttons, logging of deleted messages, or more commands by extending handlers/handlers.go.

Contributing
Feel free to open issues or submit pull requests for new features or bug fixes.
License
MIT
