# Planner Telegram Bot (Go)

Minimal Telegram "planner" bot for managing tasks with reminders.

## Quickstart

1. Create a Telegram bot token with BotFather.
2. Export env vars and run:

```bash
export BOT_TOKEN=YOUR_TOKEN
export DEFAULT_TZ=UTC

go run ./cmd/bot
```

## Commands

- `/start` - quick help
- `/help` - full help
- `/add <text>` - create a task
- `/list` - list tasks (includes inline Done/Delete buttons)
- `/done <id>` - mark done
- `/undone <id>` - mark not done
- `/del <id>` - delete task
- `/remind <id> <YYYY-MM-DD HH:MM>` - set reminder
- `/today` - list tasks with reminders today

## Notes

- Reminders use `DEFAULT_TZ` for parsing.
- Storage is in-memory only (bot restart clears tasks/reminders).

## Tests

```bash
go test ./...
```
