package main

import (
	"log"
	"net/http"

	"github.com/jon-at-github/hello-api/config"
	"github.com/jon-at-github/hello-api/handlers"
	"github.com/jon-at-github/hello-api/handlers/rest"
	"github.com/jon-at-github/hello-api/translation"
)

func main() {
	cfg := config.LoadConfiguration()
	addr := cfg.Port

	mux := http.NewServeMux()

	var translationService rest.Translator
	translationService = translation.NewStaticService()
	if cfg.LegacyEndpoint != "" {
		log.Printf("creating external translation client: %s", cfg.LegacyEndpoint)
		client := translation.NewHelloClient(cfg.LegacyEndpoint)
		translationService = translation.NewRemoteService(client)
	}

	translateHandler := rest.NewTranslateHandler(translationService)
	mux.HandleFunc("/hello", translateHandler.TranslateHandler)
	mux.HandleFunc("/health", handlers.HealthCheck)

	log.Printf("listeting on %s\n", addr)

	log.Fatal(http.ListenAndServe(addr, mux))
}

type Resp struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}
