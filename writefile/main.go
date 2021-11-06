package main

import (
	"fmt"
	"io/ioutil"
)

type Page struct {
	Title string
	Body  []byte
}

func main() {
	fmt.Println("Hello, world.")
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, err := loadPage("TestPage")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(p2.Body))
}

func (p *Page) save() error {
	fmt.Println("Saving page", p.Title)
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	fmt.Println("Loading page", title)
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}
