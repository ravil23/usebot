![CD](https://github.com/ravil23/usebot/workflows/CD/badge.svg?branch=master)

# usebot
Telegram bots for USE in Russia

## Requirements
- `Go v1.13`
- `Docker Compose` (optional)
- `Heroku CLI` (optional)

## Running
Specify `BOT_TOKEN` environment variable and run next command:
```
docker-compose up --build -d
```

## Heroku
Login one time on a host before starting work:
```
heroku login
```

Start Telegram bot:
```
heroku ps:scale -a use-bot telegrambot=1
```

Stop Telegram bot:
```
heroku ps:stop -a use-bot telegrambot
```
