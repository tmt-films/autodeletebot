<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Telegram Auto-Delete Bot</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f4f4f4;
        }
        h1, h2, h3 {
            color: #333;
        }
        h1 {
            border-bottom: 2px solid #333;
            padding-bottom: 10px;
        }
        h2 {
            margin-top: 20px;
        }
        pre, code {
            background-color: #e0e0e0;
            padding: 10px;
            border-radius: 5px;
            font-family: 'Courier New', Courier, monospace;
        }
        pre {
            overflow-x: auto;
        }
        ul {
            list-style-type: disc;
            margin-left: 20px;
        }
        a {
            color: #0066cc;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        .section {
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <h1>Telegram Auto-Delete Bot</h1>
    <p>A Telegram bot written in Go that automatically deletes messages in group chats based on time, content, user, or media type. It uses MongoDB for storing chat settings and message data, with commands like <code>/start</code> and <code>/help</code> for user interaction. Ideal for keeping group chats clean and enforcing community rules.</p>

    <div class="section">
        <h2>Features</h2>
        <ul>
            <li><strong>Time-Based Deletion</strong>: Deletes messages after a set time (e.g., <code>/settimer 5m</code> for 5 minutes).</li>
            <li><strong>Content-Based Deletion</strong>: Deletes messages containing banned words (e.g., <code>/banword spam</code>).</li>
            <li><strong>User-Specific Deletion</strong>: Deletes messages from non-admins or specific users.</li>
            <li><strong>Media-Based Deletion</strong>: Deletes photos, videos, or stickers (e.g., <code>/nomedia</code>).</li>
            <li><strong>Admin Controls</strong>: Restricts configuration commands (<code>/settimer</code>, <code>/banword</code>, <code>/nomedia</code>) to group admins.</li>
            <li><strong>Start Command</strong>: Initializes the bot and chat settings with <code>/start</code>.</li>
            <li><strong>Help Command</strong>: Lists all available commands with <code>/help</code>.</li>
            <li><strong>MongoDB Storage</strong>: Stores settings and messages with a TTL index for automatic deletion.</li>
        </ul>
    </div>

    <div class="section">
        <h2>Prerequisites</h2>
        <ul>
            <li>Go 1.21 or higher</li>
            <li>A Telegram bot token from <a href="https://t.me/BotFather">BotFather</a></li>
            <li>MongoDB instance (e.g., <a href="https://www.mongodb.com/cloud/atlas">MongoDB Atlas</a>, local MongoDB, or cloud provider)</li>
            <li>MongoDB URI and database name</li>
        </ul>
    </div>

    <div class="section">
        <h2>Setup</h2>
        <ol>
            <li>
                <strong>Clone the Repository</strong>:
                <pre><code>git clone https://github.com/YourUsername/telegram-auto-delete-bot.git
cd telegram-auto-delete-bot</code></pre>
            </li>
            <li>
                <strong>Install Dependencies</strong>:
                <pre><code>go mod tidy</code></pre>
            </li>
            <li>
                <strong>Set Environment Variables</strong>:
                <p>Create a <code>.env</code> file or set the following:</p>
                <pre><code>export BOT_TOKEN="your-bot-token-here"
export MONGO_URI="mongodb+srv://&lt;username&gt;:&lt;password&gt;@&lt;cluster&gt;.mongodb.net/?retryWrites=true&w=majority"
export MONGO_DB_NAME="telegram_bot"</code></pre>
            </li>
            <li>
                <strong>Run the Bot Locally</strong>:
                <pre><code>go run main.go</code></pre>
            </li>
            <li>
                <strong>Add the Bot to a Group</strong>:
                <ul>
                    <li>Add the bot to your Telegram group and make it an admin with "Delete Messages" permission.</li>
                    <li>Use <code>/setcommands</code> in BotFather to register commands:</li>
                </ul>
                <pre><code>start - Initialize the bot
help - Show available commands
settimer - Set auto-delete timer (e.g., /settimer 5m)
banword - Ban a word (e.g., /banword spam)
nomedia - Delete all media messages</code></pre>
            </li>
        </ol>
    </div>

    <div class="section">
        <h2>Usage</h2>
        <ul>
            <li><code>/start</code>: Initializes the bot and creates default chat settings. Responds with a welcome message.</li>
            <li><code>/help</code>: Lists all commands with descriptions.</li>
            <li><code>/settimer &lt;minutes&gt;m</code>: Sets a timer to delete messages after the specified time (e.g., <code>/settimer 5m</code>) [Admin only].</li>
            <li><code>/banword &lt;word&gt;</code>: Deletes messages containing the specified word (e.g., <code>/banword spam</code>) [Admin only].</li>
            <li><code>/nomedia</code>: Deletes all photos, videos, and stickers [Admin only].</li>
        </ul>
    </div>

    <div class="section">
        <h2>MongoDB Setup</h2>
        <ol>
            <li>Create a MongoDB database (e.g., via <a href="https://www.mongodb.com/cloud/atlas">MongoDB Atlas</a> free tier or local MongoDB with <code>docker run -d -p 27017:27017 mongo</code>).</li>
            <li>Get the connection URI (e.g., <code>mongodb+srv://&lt;username&gt;:&lt;password&gt;@&lt;cluster&gt;.mongodb.net/?retryWrites=true&w=majority</code>).</li>
            <li>Set <code>MONGO_URI</code> and <code>MONGO_DB_NAME</code> in your environment.</li>
            <li>The bot uses two collections:
                <ul>
                    <li><code>settings</code>: Stores chat-specific settings (<code>chat_id</code>, <code>timer</code>, <code>no_media</code>, <code>banned_words</code>).</li>
                    <li><code>messages</code>: Stores messages with a TTL index for auto-deletion based on the <code>expireAt</code> field.</li>
                </ul>
            </li>
        </ol>
    </div>

    <div class="section">
        <h2>Deployment</h2>
        <p>To run the bot 24/7, deploy on a cloud platform like Heroku, AWS, or Render:</p>
        <ol>
            <li>Set up a Go environment.</li>
            <li>Configure <code>BOT_TOKEN</code>, <code>MONGO_URI</code>, and <code>MONGO_DB_NAME</code> as environment variables.</li>
            <li>Ensure MongoDB is accessible (e.g., whitelist your server’s IP in MongoDB Atlas).</li>
            <li>Use a persistent MongoDB instance (e.g., MongoDB Atlas for cloud hosting).</li>
            <li>Example for Heroku:
                <pre><code>heroku create
heroku config:set BOT_TOKEN="your-bot-token" MONGO_URI="your-mongo-uri" MONGO_DB_NAME="telegram_bot"
git push heroku main</code></pre>
            </li>
        </ol>
    </div>

    <div class="section">
        <h2>Testing</h2>
        <ol>
            <li>Add the bot to a private Telegram group and make it an admin.</li>
            <li>Test commands:
                <ul>
                    <li><code>/start</code>: Should respond with a welcome message and create a <code>settings</code> document in MongoDB.</li>
                    <li><code>/help</code>: Should list all commands.</li>
                    <li><code>/settimer 1m</code>: Set a 1-minute deletion timer (admin only).</li>
                    <li><code>/banword test</code>: Add "test" as a banned word (admin only).</li>
                    <li><code>/nomedia</code>: Enable media deletion (admin only).</li>
                    <li>Send messages with banned words, media (e.g., a photo), or wait for the timer to verify deletions.</li>
                </ul>
            </li>
            <li>Check MongoDB collections (<code>settings</code> and <code>messages</code>) to confirm data storage.</li>
        </ol>
    </div>

    <div class="section">
        <h2>Notes</h2>
        <ul>
            <li><strong>Replace <code>YourUsername</code></strong>: Update <code>go.mod</code> and <code>README.html</code> with your GitHub username.</li>
            <li><strong>Error Handling</strong>: The bot logs errors to the console. For production, consider adding a logging service (e.g., to a file or external service).</li>
            <li><strong>Security</strong>: Admin-only commands are restricted using <code>isAdmin</code>. Ensure the bot has "Delete Messages" permission in groups.</li>
            <li><strong>Scalability</strong>: MongoDB’s TTL index handles time-based deletions efficiently. For high-traffic groups, consider sharding or replica sets in MongoDB.</li>
            <li><strong>Enhancements</strong>: Add inline buttons, logging of deleted messages, or more commands by extending <code>handlers/handlers.go</code>.</li>
        </ul>
    </div>

    <div class="section">
        <h2>Contributing</h2>
        <p>Feel free to open issues or submit pull requests for new features or bug fixes.</p>
    </div>

    <div class="section">
        <h2>License</h2>
        <p>MIT</p>
    </div>
</body>
</html>
