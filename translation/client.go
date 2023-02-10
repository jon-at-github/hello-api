package translation

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

var _ HelloClient = &APIClient{}

// APIClient will store a given endpoint for external calls.
type APIClient struct {
	endpoint string
}

// NewHelloClient creates instance of client with a given endpoint.
func NewHelloClient(endpoint string) *APIClient {
	return &APIClient{
		endpoint: endpoint,
	}
}

// Translate will call external client for translation.
func (c *APIClient) Translate(word, language string) (string, error) {
	request := map[string]interface{}{
		"word":     word,
		"language": language,
	}
	b, err := json.Marshal(request)
	if err != nil {
		return "", errors.New("unable to encode message")
	}
	response, err := http.Post(c.endpoint, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Println(err)
		return "", errors.New("call to api failed")
	}

	if response.StatusCode == http.StatusNotFound {
		return "", nil
	}

	if response.StatusCode == http.StatusInternalServerError {
		return "", errors.New("error in api")
	}

	b, _ = io.ReadAll(response.Body)
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return "", errors.New("unable to decode message")
	}

	return m["translation"].(string), nil
}
