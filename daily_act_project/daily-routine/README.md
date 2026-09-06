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

## Deploying to Fly.io

This app is ready to deploy — it reads `PORT` and `DB_PATH` from the environment (both default sensibly for local dev).

1. Install the CLI: `curl -L https://fly.io/install.sh | sh` (Mac/Linux) or `iwr https://fly.io/install.ps1 -useb | iex` (Windows)
2. `fly auth signup` (or `fly auth login`)
3. From this folder, run `fly launch` — it will detect the Dockerfile, ask for an app name (becomes `your-app-name.fly.dev`) and a region. Say **no** to a Postgres database.
4. Create a persistent volume so your data survives redeploys:
   `fly volumes create routine_data --size 1`
   Then add this to the `fly.toml` file `fly launch` generated:
   ```toml
   [mounts]
     source = "routine_data"
     destination = "/data"
   ```
5. `fly deploy`
6. Your site is live at `https://your-app-name.fly.dev`

Redeploy any time with `fly deploy` from this folder.

## Possible next steps


- Recurring activities (a daily template you copy into each new day)
- Streaks / stats page
- Edit an existing activity (currently you'd delete + re-add)
- Basic auth if you want to deploy this somewhere shared
