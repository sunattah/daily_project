# Daily Routine Tracker

A simple Go web app for logging your daily activities by time and date, backed by SQLite.

## Run it

```bash
go mod tidy
go run .
```

Then open http://localhost:8080

## Features

- View activities for any date (date picker at the top)
- Add an activity with a date, time, and description (`/add`)
- Mark an activity complete/incomplete (✔ button)
- Delete an activity (✕ button)

## Project structure

- `main.go` — server setup, DB access (SQLite via `database/sql`), and route handlers
- `templates/index.html` — the day's activity list
- `templates/add.html` — the "add activity" form
- `routine.db` — created automatically on first run (SQLite database file)

## Possible next steps

- Recurring activities (a daily template you copy into each new day)
- Streaks / stats page
- Edit an existing activity (currently you'd delete + re-add)
- Basic auth if you want to deploy this somewhere shared
