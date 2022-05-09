package authorization

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthorizationTestSuite struct {
	suite.Suite
}

func (suite *AuthorizationTestSuite) TestAuthorization_AuthorizeAdmin() {
	userId := "test-id"

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/", nil)

	ctx.Request.Header.Set("X-User-Id", userId)
	ctx.Request.Header.Set("X-User-Claims", `{"admin": true}`)

	sut := NewRest(ctx)

	suite.Equal(userId, sut.id)

	suite.True(sut.AuthorizeAdmin())
}

func (suite *AuthorizationTestSuite) TestAuthorization_AuthorizeAdmin_NoAdmin() {
	userId := "test-id"

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/", nil)

	ctx.Request.Header.Set("X-User-Id", userId)
	ctx.Request.Header.Set("X-User-Claims", `{"admin": false}`)

	sut := NewRest(ctx)

	suite.Equal(userId, sut.id)

	suite.False(sut.AuthorizeAdmin())
}

func (suite *AuthorizationTestSuite) TestAuthorization_AuthorizeAdmin_NoClaims() {
	userId := "test-id"

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/", nil)

	ctx.Request.Header.Set("X-User-Id", userId)

	sut := NewRest(ctx)

	suite.Equal(userId, sut.id)

	suite.False(sut.AuthorizeAdmin())
}

func (suite *AuthorizationTestSuite) TestAuthorization_AuthorizeMatchingId() {
	userId := "test-id"

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/", nil)

	ctx.Request.Header.Set("X-User-Id", userId)

	sut := NewRest(ctx)

	suite.Equal(userId, sut.id)

	suite.True(sut.AuthorizeMatchingId(userId))
}

func (suite *AuthorizationTestSuite) TestAuthorization_AuthorizeMatchingId_NoMatch() {
	userId := "test-id"
	expected := "test-2"

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request, _ = http.NewRequest("GET", "/", nil)

	ctx.Request.Header.Set("X-User-Id", userId)

	sut := NewRest(ctx)

	suite.Equal(userId, sut.id)

	suite.False(sut.AuthorizeMatchingId(expected))
}

func TestUnit_AuthorizationTestSuite(t *testing.T) {
	repoSuite := new(AuthorizationTestSuite)
	suite.Run(t, repoSuite)
}
