package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
	"text/template"
)

type Page struct {
	Title string
	Body  []byte
}

var List = []string{}

// Caching Templates
var templates = template.Must(template.ParseFiles("./templ/edit.html", "./templ/view.html", "./templ/list.html"))

// Filter Path
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// File Handler: Save
func (p *Page) save() error {
	fmt.Println("Saving page", p.Title)
	filename := "./data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

// File Handler: Load
func loadPage(title string) (*Page, error) {
	fmt.Println("Loading page", title)
	filename := "./data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// Get sanitised titles
// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
// 	m := validPath.FindStringSubmatch(r.URL.Path)
// 	if m == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("invalid Page Title")
// 	}
// 	return m[2], nil // The title is the second subexpression.
// }

func makeHandeler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// Rendering Templates
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	template := tmpl + ".html"
	err := templates.ExecuteTemplate(w, template, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler function to view a page
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// Handler function to edit a page
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// Handler function to save a page
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// List all saved pages
func listHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir("./data/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(files) != len(List) {
		for _, f := range files {
			title := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
			// fmt.Fprintf(w, "%s\n", title)
			List = append(List, title)
		}
	}

	// fmt.Println(List)
	// fmt.Println(List[0])
	err = templates.ExecuteTemplate(w, "list.html", List)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Main function
func main() {
	// p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	// p1.save()
	http.HandleFunc("/view/", makeHandeler(viewHandler))
	http.HandleFunc("/edit/", makeHandeler(editHandler))
	http.HandleFunc("/save/", makeHandeler(saveHandler))
	http.HandleFunc("/list/", listHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
