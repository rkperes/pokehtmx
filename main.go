package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/rkperes/pokehtmx/templates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	component := templates.Index()

	mux := http.NewServeMux()

	mux.Handle("/", templ.Handler(component))
	mux.HandleFunc("POST /search", searchHandler)

	fmt.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalf("http.ListenAndServe: %v", err)
	}
}

const pokeapiSearchPokemon = "https://pokeapi.co/api/v2/pokemon"

type pokepapiSearchPokemonResponse struct {
	Id      int            `json:"id"`
	Name    string         `json:"name"`
	Sprites pokeapiSprites `json:"sprites"`
}

type pokeapiSprites struct {
	FrontDefault string `json:"front_default"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := strings.ToLower(r.FormValue("pokemon"))
	if len(searchQuery) == 0 {
		slog.Error("pokeapi search empty", slog.Any("form", r.Form))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := fmt.Sprintf("%s/%s", pokeapiSearchPokemon, searchQuery)
	resp, err := http.Get(u)
	if err != nil {
		slog.Error(fmt.Sprintf("pokeapi search: %v", err), slog.String("url", u))
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		slog.Warn("pokeapi not found", slog.String("url", u))
		http.Error(w, "Pokemon not found", http.StatusNotFound)
		return
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error(fmt.Sprintf("pokeapi error: %v", err), slog.String("url", u))
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	var pokeapiResp pokepapiSearchPokemonResponse
	if err := json.NewDecoder(resp.Body).Decode(&pokeapiResp); err != nil {
		slog.Error(fmt.Sprintf("pokeapi unmarshal: %v", err), slog.String("url", u))
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	component := templates.SearchResultDisplay(templates.SearchResult{
		Id:        fmt.Sprintf("%03d", pokeapiResp.Id),
		Name:      cases.Title(language.English).String(pokeapiResp.Name),
		SpriteURL: pokeapiResp.Sprites.FrontDefault,
	})

	slog.Info("serve search results", slog.Any("pokemon", pokeapiResp), slog.String("url", u))
	templ.Handler(component).ServeHTTP(w, r)
}
