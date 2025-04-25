package core

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var testToken = "testtokentesttoken"
var testSub = "testID"

func TestAuthMiddleware(t *testing.T) {
	response := httptest.NewRecorder()
	c, engine := gin.CreateTestContext(response)

	c.Request, _ = http.NewRequest(http.MethodGet, "/deviceList", nil)
	c.Request.Header.Add("Authorization", testToken)

	mocked := &dataFetcherMock{}
	mockObj := mocked.On("GetUserInfo", context.Background(), testToken, "").
		Return(&gocloak.UserInfo{Sub: &testSub}, nil)
	defer func() { mockObj.Unset() }()
	engine.Use(AuthMiddleware(mocked))
	engine.GET("/deviceList", func(c *gin.Context) {
		c.JSON(http.StatusOK, "success")
	})
	engine.ServeHTTP(response, c.Request)
	require.Equal(t, http.StatusOK, response.Result().StatusCode)
}

type dataFetcherMock struct {
	mock.Mock
}

func (k *dataFetcherMock) GetUserInfo(ctx context.Context, token string, realm string) (*gocloak.UserInfo, error) {
	args := k.Called(ctx, token, realm)
	us := args.Get(0)
	if us == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gocloak.UserInfo), nil
}
