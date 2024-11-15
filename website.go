package main

import (
	"html/template"
	"net/http"
	"fmt"
)

func ServeWebsite(dal DataAccessLayer) {
	tmpl := template.Must(template.ParseFiles("submission_form.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		var completeness bool
		if r.FormValue("complete") == "true" {
			completeness = true
		} else if r.FormValue("complete") == "false" {
			completeness = false
		} else {
			fmt.Println("incorrect form value")
		}
		todo := ToDoItem{
			Id:       Id(r.FormValue("id")),
			Title:    Title(r.FormValue("title")),
			Priority: Priority(r.FormValue("priority")),
			Complete: Complete(completeness),
		}

		err := dal.db.create(todo)

		if err != nil {
			fmt.Printf("ERROR!: %q", err)
			tmpl.Execute(w, struct{ Success bool }{false})
		} else {
			tmpl.Execute(w, struct{ Success bool }{true})
		}

	})

	http.ListenAndServe(":6060", nil)
}
