package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/rkperes/pokehtmx/templates"
)

func main() {
	component := templates.Index()

	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}
