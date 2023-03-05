package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// the handler to display our page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	// start the web server
	fmt.Println("Starting front end service on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFS embed.FS

// render generates a page of html from our template files
func render(w http.ResponseWriter, t string) {
        tmpl,err:=template.ParseGlob("templates/*.gohtml")
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        // execute the template
        if err := tmpl.ExecuteTemplate(w,t,nil); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}
