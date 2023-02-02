// Package faas is used for function definitions
package faas

import (
	"net/http"

	"github.com/jon-at-github/hello-api/handlers/rest"
	"github.com/jon-at-github/hello-api/translation"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	translationService := translation.NewStaticService()
	translateHandler := rest.NewTranslateHandler(translationService)
	translateHandler.TranslateHandler(w, r)
}
