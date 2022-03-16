package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type TemplateData struct {
	Data            map[string]interface{}
	IsAuthenticated int
}

func (app *Config) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	if app.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = 1
	}
	return td
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	log.Println("rendering template", t)
	partials := []string{
		"./templates/base.layout.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if td == nil {
		td = &TemplateData{}
	}
	td = app.addDefaultData(td, r)

	if err := tmpl.Execute(w, td); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
