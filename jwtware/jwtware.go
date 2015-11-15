package jwtware

import (
	"io/ioutil"
	"net/http"

	"github.com/InteractiveLecture/serviceclient"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var secret = make([]byte, 0)

func SecretHandler(token *jwt.Token) (interface{}, error) {
	if len(secret) == 0 {
		serviceInstance := serviceclient.GetInstance("authentication-service")
		resp, err := serviceInstance.Get("/oauth/token_key")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, decodeErr := ioutil.ReadAll(resp.Body)
		if decodeErr != nil {
			return nil, err
		}
		secret = body
	}
	return secret, nil
}

func New(next http.Handler) http.Handler {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: SecretHandler,
	})
	return jwtMiddleware.Handler(next)
}
