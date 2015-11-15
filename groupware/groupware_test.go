package groupware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/InteractiveLecture/middlewares/groupware/mocks"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func prepareMockToken(t *testing.T, authorities ...string) (*gomock.Controller, *jwt.Token) {
	controller := gomock.NewController(t)
	mockMethod := mocks.NewMockSigningMethod(controller)
	mockMethod.EXPECT().Alg().Return("mock")
	token := jwt.New(mockMethod)
	token.Claims = map[string]interface{}{
		"authorities": authorities,
	}
	return controller, token
}

func TestNew(t *testing.T) {
	controller, token := prepareMockToken(t, "admin")
	defer controller.Finish()

	server := httptest.NewServer(createHandlerChain(token, "admin", "officer"))
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.Nil(t, err)
	// if statuscode == teapot, the last handler was called.
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)

}

func TestFail(t *testing.T) {
	controller, token := prepareMockToken(t, "user")
	defer controller.Finish()
	server := httptest.NewServer(createHandlerChain(token, "admin", "officer"))
	resp, err := http.Get(server.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestMultipleGroups(t *testing.T) {
	controller, token := prepareMockToken(t, "user", "officer", "assistant")
	defer controller.Finish()
	server := httptest.NewServer(createHandlerChain(token, "officer"))
	resp, err := http.Get(server.URL)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)

}

func createHandlerChain(token *jwt.Token, permissions ...string) http.Handler {
	checkHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	return createContextHandler(New(DefaultOptions(checkHandler, permissions...)), token)
}

func createContextHandler(next http.Handler, token *jwt.Token) http.Handler {
	contextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "user", token)
		next.ServeHTTP(w, r)
	})
	return contextHandler
}
