package translation_test

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jon-at-github/hello-api/translation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HelloClientSuite struct {
	suite.Suite
	mockServerService *MockService
	server            *httptest.Server
	underTest         translation.HelloClient
}

func TestHelloClientSuite(t *testing.T) {
	suite.Run(t, new(HelloClientSuite))
}

type MockService struct {
	mock.Mock
}

func (m *MockService) Translate(word, language string) (string, error) {
	args := m.Called(word, language)
	return args.String(0), args.Error(1)
}

func (suite *HelloClientSuite) SetupSuite() {
	suite.mockServerService = new(MockService)
	handler := func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		defer func() {
			err := r.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		var m map[string]interface{}
		_ = json.Unmarshal(b, &m)

		word := m["word"].(string)
		language := m["language"].(string)

		response, err := suite.mockServerService.Translate(word, language)
		if err != nil {
			http.Error(w, "error", 500)
		}
		if response == "" {
			http.Error(w, "missing", 404)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, response)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	suite.server = httptest.NewServer(mux)
	suite.underTest = translation.NewHelloClient(suite.server.URL)
}

func (suite *HelloClientSuite) SetupTest() {
	suite.mockServerService = new(MockService)
}

func (suite *HelloClientSuite) TeardownSuite() {
	suite.server.Close()
}

func (suite *HelloClientSuite) TestCall() {
	suite.mockServerService.On("Translate", "foo", "bar").Return(`{"translation":"baz"}`, nil)

	response, err := suite.underTest.Translate("foo", "bar")

	suite.NoError(err)
	suite.Equal(response, "baz")
}

func (suite *HelloClientSuite) TestCall_APIError() {
	suite.mockServerService.On("Translate", "foo", "bar").Return("", errors.New("this is a test"))

	response, err := suite.underTest.Translate("foo", "bar")

	suite.EqualError(err, "error in api")
	suite.Equal(response, "")
}

func (suite *HelloClientSuite) TestCall_InvalidJSON() {
	suite.mockServerService.On("Translate", "foo", "bar").Return(`invalid json`, nil)

	response, err := suite.underTest.Translate("foo", "bar")

	suite.EqualError(err, "unable to decode message")
	suite.Equal(response, "")
}
