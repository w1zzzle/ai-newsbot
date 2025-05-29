# AI NewsBot

A self-contained Go service that scrapes Reddit posts, translates them to Russian using AI, and publishes them to Telegram.

## Features

- 🔍 **Reddit Scraping**: Extracts hot posts from configured subreddits
- 🤖 **AI Translation**: Uses OpenRouter DeepSeek R1 for Russian translation
- 📱 **Telegram Publishing**: Automatically publishes to Telegram channels
- 🗄️ **PostgreSQL Storage**: Persistent storage with duplicate detection
- ⏰ **Scheduled Execution**: Runs hourly via cron
- 🐳 **Docker Support**: Full containerization with docker-compose
- 🧪 **Comprehensive Testing**: Unit tests with mocking
- 🚀 **CI/CD Ready**: GitHub Actions workflow included

## Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Scraper   │───▶│  Storage    │───▶│ Translation │───▶│    Bot      │
│             │    │             │    │             │    │             │
│ Reddit HTML │    │ PostgreSQL  │    │ OpenRouter  │    │  Telegram   │
└─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘
```

## Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/youruser/ai-newsbot.git
   cd ai-newsbot
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your API keys and configuration
   ```

3. **Run with Docker Compose**
   ```bash
   docker-compose up -d
   ```

4. **Or run locally**
   ```bash
   go mod download
   go run cmd/ai-newsbot/main.go
   ```

## Configuration

| Environment Variable | Description | Required |
|---------------------|-------------|----------|
| `POSTGRES_DSN` | PostgreSQL connection string | Yes |
| `REDDIT_URLS` | Comma-separated Reddit URLs | No |
| `UPVOTE_THRESHOLD` | Minimum upvotes for posts | No |
| `OPENROUTER_API_KEY` | OpenRouter API key | Yes |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token | Yes |
| `TELEGRAM_CHAT_ID` | Telegram chat/channel ID | Yes |

## Development

### Prerequisites
- Go 1.21+
- PostgreSQL 12+
- Docker & Docker Compose (optional)

### Running Tests
```bash
go test -v ./...
go test -v -race -coverprofile=coverage.out ./...
```

### Database Setup
```bash
# Create database
createdb ai_newsbot

# Run migrations
psql ai_newsbot < migrations/schema.sql
```

### Linting
```bash
golangci-lint run ./...
```

## API Keys Setup

1. **OpenRouter**: Get your API key from [OpenRouter](https://openrouter.ai/)
2. **Telegram Bot**: Create a bot via [@BotFather](https://t.me/botfather)
3. **Telegram Chat ID**: Use [@userinfobot](https://t.me/userinfobot) to get your chat ID

## Deployment

### Docker
```bash
docker build -t ai-newsbot .
docker run -d --env-file .env ai-newsbot
```

### Docker Compose
```bash
docker-compose up -d
```

## Project Structure

```
ai-newsbot/
├── cmd/ai-newsbot/          # Application entry point
├── internal/
│   ├── app/                 # Application logic
│   ├── bot/                 # Telegram bot integration
│   ├── config/              # Configuration management
│   ├── scraper/             # Reddit scraping logic
│   ├── storage/             # Database operations
│   └── translation/         # AI translation service
├── migrations/              # Database schemas
├── .github/workflows/       # CI/CD pipelines
└── docker-compose.yml       # Container orchestration