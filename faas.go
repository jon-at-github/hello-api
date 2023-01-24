// Package faas is used for function definitions
package faas

import (
	"net/http"

	"github.com/jon-at-github/hello-api/handlers/rest"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	rest.TranslateHandler(w, r)
}
