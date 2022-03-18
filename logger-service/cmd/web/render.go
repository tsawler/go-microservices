package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type TemplateData struct {
	Data            map[string]any
	IsAuthenticated int
}

// addDefaultData adds whatever is specified in this function to all templates as
// data that can be accessed directly.
func (app *Config) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	if app.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = 1
	}
	return td
}

func (app *Config) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	// we only have one partial, which is actually a layout, but since all of our pages
	// require the layout, we need to include it when we call ParseFiles, below.
	// If you have other partials you use in your templates, add them to this slice.
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
