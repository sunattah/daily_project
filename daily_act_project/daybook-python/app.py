import os
import sqlite3
from datetime import date

from flask import Flask, g, redirect, render_template, request, url_for

app = Flask(__name__)

DB_PATH = os.environ.get("DB_PATH", "routine.db")


def get_db():
    """Open (or reuse) a per-request SQLite connection."""
    if "db" not in g:
        g.db = sqlite3.connect(DB_PATH)
        g.db.row_factory = sqlite3.Row
    return g.db


@app.teardown_appcontext
def close_db(exception=None):
    db = g.pop("db", None)
    if db is not None:
        db.close()


def init_db():
    db = sqlite3.connect(DB_PATH)
    db.execute(
        """
        CREATE TABLE IF NOT EXISTS activities (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            time TEXT NOT NULL,
            description TEXT NOT NULL,
            completed BOOLEAN NOT NULL DEFAULT 0
        )
        """
    )
    db.commit()
    db.close()


def get_activities(for_date):
    db = get_db()
    rows = db.execute(
        "SELECT id, date, time, description, completed "
        "FROM activities WHERE date = ? ORDER BY time ASC",
        (for_date,),
    ).fetchall()
    return rows


def create_activity(for_date, time_str, description):
    db = get_db()
    db.execute(
        "INSERT INTO activities (date, time, description, completed) "
        "VALUES (?, ?, ?, 0)",
        (for_date, time_str, description),
    )
    db.commit()


def toggle_activity(activity_id):
    db = get_db()
    db.execute(
        "UPDATE activities SET completed = NOT completed WHERE id = ?",
        (activity_id,),
    )
    db.commit()


def delete_activity(activity_id):
    db = get_db()
    db.execute("DELETE FROM activities WHERE id = ?", (activity_id,))
    db.commit()


@app.route("/")
def index():
    for_date = request.args.get("date") or date.today().isoformat()
    activities = get_activities(for_date)
    return render_template("index.html", date=for_date, activities=activities)


@app.route("/add", methods=["GET", "POST"])
def add():
    if request.method == "POST":
        for_date = request.form.get("date", "")
        time_str = request.form.get("time", "")
        description = request.form.get("description", "")

        if not (for_date and time_str and description):
            return "date, time, and description are required", 400

        create_activity(for_date, time_str, description)
        return redirect(url_for("index", date=for_date))

    for_date = request.args.get("date") or date.today().isoformat()
    return render_template("add.html", date=for_date)


@app.route("/toggle", methods=["POST"])
def toggle():
    activity_id = request.args.get("id", type=int)
    for_date = request.args.get("date", "")
    if activity_id is None:
        return "invalid id", 400
    toggle_activity(activity_id)
    return redirect(url_for("index", date=for_date))


@app.route("/delete", methods=["POST"])
def delete():
    activity_id = request.args.get("id", type=int)
    for_date = request.args.get("date", "")
    if activity_id is None:
        return "invalid id", 400
    delete_activity(activity_id)
    return redirect(url_for("index", date=for_date))


if __name__ == "__main__":
    init_db()
    port = int(os.environ.get("PORT", 8080))
    app.run(host="0.0.0.0", port=port)
# f76ffytfytutfyt