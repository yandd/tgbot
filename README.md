# tgbot
telegram bot: rss, ...

## How to run

1. Ceate config file `tgbot.json`
```json
{
  "telebot_token": "<YOUR BOT TOKEN>",
  "fetch_interval": 300,
  "sqlite_dbfile": "tgbot.db"
}
```

2. Run
```bash
./tgbot --config=tgbot.json
```

## Feature

* [x] Rss
