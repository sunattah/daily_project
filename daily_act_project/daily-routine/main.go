package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Activity struct {
	ID          int
	Date        string
	Time        string
	Description string
	Completed   bool
}

var db *sql.DB
var tmpl *template.Template

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "routine.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS activities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		time TEXT NOT NULL,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0
	);`

	if _, err = db.Exec(createTable); err != nil {
		log.Fatal(err)
	}
}

// getActivities returns all activities for a given date, ordered by time.
func getActivities(date string) ([]Activity, error) {
	rows, err := db.Query(
		`SELECT id, date, time, description, completed
		 FROM activities WHERE date = ? ORDER BY time ASC`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var a Activity
		if err := rows.Scan(&a.ID, &a.Date, &a.Time, &a.Description, &a.Completed); err != nil {
			return nil, err
		}
		activities = append(activities, a)
	}
	return activities, rows.Err()
}

func createActivity(date, timeStr, description string) error {
	_, err := db.Exec(
		`INSERT INTO activities (date, time, description, completed)
		 VALUES (?, ?, ?, 0)`, date, timeStr, description)
	return err
}

func toggleActivity(id int) error {
	_, err := db.Exec(
		`UPDATE activities SET completed = NOT completed WHERE id = ?`, id)
	return err
}

func deleteActivity(id int) error {
	_, err := db.Exec(`DELETE FROM activities WHERE id = ?`, id)
	return err
}

// indexHandler shows the activities for a given date (defaults to today).
func indexHandler(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	activities, err := getActivities(date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Date       string
		Activities []Activity
	}{Date: date, Activities: activities}

	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		date := r.FormValue("date")
		timeStr := r.FormValue("time")
		description := r.FormValue("description")

		if date == "" || timeStr == "" || description == "" {
			http.Error(w, "date, time, and description are required", http.StatusBadRequest)
			return
		}

		if err := createActivity(date, timeStr, description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/?date="+date, http.StatusSeeOther)
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	data := struct{ Date string }{Date: date}
	if err := tmpl.ExecuteTemplate(w, "add.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	date := r.URL.Query().Get("date")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := toggleActivity(idInt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/?date="+date, http.StatusSeeOther)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	date := r.URL.Query().Get("date")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := deleteActivity(idInt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/?date="+date, http.StatusSeeOther)
}

func main() {
	initDB()
	defer db.Close()

	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/add", addHandler)
	mux.HandleFunc("/toggle", toggleHandler)
	mux.HandleFunc("/delete", deleteHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
