package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Task struct {
	ID        int
	Title     string
	Priority  string
	Completed bool
}
type Activity struct {
	ID       int
	Time     string
	Activity string
}

var schedule []Activity

var nextActivityID = 1

var tasks []Task
var nextID = 1

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		return
	}

	data := struct {
		Tasks []Task
	}{
		Tasks: tasks,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Unable to display page", http.StatusInternalServerError)
		return
	}
}
func showSchedule(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/schedule.html")

	if err != nil {
		http.Error(w, "Unable to load schedule", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, "Unable to display schedule", http.StatusInternalServerError)
		return
	}
}

func addTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	priority := r.FormValue("priority")

	if title == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	task := Task{
		ID:        nextID,
		Title:     title,
		Priority:  priority,
		Completed: false,
	}

	tasks = append(tasks, task)
	nextID++

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func completeTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := 0

	_, err := fmt.Sscanf(r.FormValue("id"), "%d", &id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Completed = !tasks[i].Completed
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := 0

	_, err := fmt.Sscanf(r.FormValue("id"), "%d", &id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/add-task", addTask)
	http.HandleFunc("/complete-task", completeTask)
	http.HandleFunc("/delete-task", deleteTask)
	http.HandleFunc("/schedule", showSchedule)

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	log.Println("Server running at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
