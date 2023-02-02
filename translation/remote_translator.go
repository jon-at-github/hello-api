package translation

import (
	"log"
	"strings"

	"github.com/jon-at-github/hello-api/handlers/rest"
)

var _ rest.Translator = &RemoteService{}

// RemoteService  will allow for external calls to existing service for translations.
type RemoteService struct {
	client HelloClient
}

// HelloClient will call external service.
type HelloClient interface {
	Translate(word, language string) (string, error)
}

// NewRemoteService creates a new implementation of RemoteService.
func NewRemoteService(client HelloClient) *RemoteService {
	return &RemoteService{client: client}
}

// Translate will take a given word and try to find the result using the client.
func (s *RemoteService) Translate(word string, language string) string {
	word = strings.ToLower(word)
	language = strings.ToLower(language)
	response, err := s.client.Translate(word, language)
	if err != nil {
		log.Println(err)
		return ""
	}
	return response
}
