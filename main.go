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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		slog.Warn(fmt.Sprintf("pokeapi not found: %v", err), slog.String("url", u))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var pokeapiResp pokepapiSearchPokemonResponse
	if err := json.NewDecoder(resp.Body).Decode(&pokeapiResp); err != nil {
		slog.Error(fmt.Sprintf("pokeapi unmarshal: %v", err), slog.String("url", u))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	component := templates.SearchResultDisplay(templates.SearchResult{
		Name:      pokeapiResp.Name,
		SpriteURL: pokeapiResp.Sprites.FrontDefault,
	})

	slog.Info("serve search results", slog.Any("pokemon", pokeapiResp), slog.String("url", u))
	templ.Handler(component).ServeHTTP(w, r)
}
