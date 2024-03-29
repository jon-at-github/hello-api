// Package rest houses all rest handlers
package rest

import (
	"encoding/json"
	"net/http"
	"strings"
)

// TranslateHandler will translate calls for caller.
type TranslateHandler struct {
	service Translator
}

type Translator interface {
	Translate(word string, language string) string
}

// NewTranslateHandler will create a new instance of the handler using a
// translation service.
func NewTranslateHandler(service Translator) *TranslateHandler {
	return &TranslateHandler{
		service: service,
	}
}

type Resp struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}

func (t *TranslateHandler) TranslateHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	language := r.URL.Query().Get("language")
	if language == "" {
		language = "english"
	}
	word := strings.ReplaceAll(r.URL.Path, "/", "")
	translation := t.service.Translate(word, language)
	if translation == "" {
		w.WriteHeader(404)
		return
	}
	resp := Resp{
		Language:    language,
		Translation: translation,
	}
	if err := enc.Encode(resp); err != nil {
		panic("unable to encode response")
	}
}
