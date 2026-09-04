package models

import (
	"daily_life/models"
	"html/template"
	"net/http"
)

type Note struct {
	ID      int
	Title   string
	Content string
}

var notes []models.Note
var nextNoteID = 1

func showNotes(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/notes.html")

	if err != nil {
		http.Error(w, "Unable to load notes", http.StatusInternalServerError)
		return
	}

	data := struct {
		Notes []Note
	}{
		Notes: notes,
	}

	err = tmpl.Execute(w, data)

	if err != nil {
		http.Error(w, "Unable to display notes", http.StatusInternalServerError)
		return
	}
}
