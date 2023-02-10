package translation_test

import (
	"errors"
	"testing"

	"github.com/jon-at-github/hello-api/translation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestRemoteServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteServiceTestSuite))
}

type RemoteServiceTestSuite struct {
	suite.Suite
	client    *MockHelloClient
	underTest *translation.RemoteService
}

func (suite *RemoteServiceTestSuite) SetupTest() {
	suite.client = new(MockHelloClient)
	suite.underTest = translation.NewRemoteService(suite.client)
}

type MockHelloClient struct {
	mock.Mock
}

func (m *MockHelloClient) Translate(word, language string) (string, error) {
	args := m.Called(word, language)
	return args.String(0), args.Error(1)
}

func (suite *RemoteServiceTestSuite) TestTranslate() {
	suite.client.On("Translate", "foo", "bar").Return("baz", nil)

	result := suite.underTest.Translate("Foo", "bar")

	suite.Equal(result, "baz")
	suite.client.AssertExpectations(suite.T())
}

func (suite *RemoteServiceTestSuite) TestTranslate_Error() {
	suite.client.On("Translate", "foo", "bar").Return("baz", errors.New("failure"))

	result := suite.underTest.Translate("foo", "bar")

	suite.Equal(result, "")
	suite.client.AssertExpectations(suite.T())
}

func (suite *RemoteServiceTestSuite) TestTranslate_Cache() {
	suite.client.On("Translate", "foo", "bar").Return("baz", nil).Times(1)

	result1 := suite.underTest.Translate("Foo", "bar")
	result2 := suite.underTest.Translate("foo", "bar")

	suite.Equal("baz", result1)
	suite.Equal("baz", result2)
	suite.client.AssertExpectations(suite.T())
}
