package jwtware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/InteractiveLecture/serviceclient"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

var secret []byte

func SecretHandler(token *jwt.Token) (interface{}, error) {
	if secret == nil {
		serviceInstance := serviceclient.New("authentication-service")
		resp, err := serviceInstance.Get("/oauth/token_key")
		if err != nil {
			log.Println("error occured while requesting key: ", err)
			return nil, err
		}
		defer resp.Body.Close()
		result := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Println("error occured while reading key: ", err)
			return nil, err
		}
		secret = []byte(result["value"].(string))
	}
	return secret, nil
}

func New(next http.Handler) http.Handler {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: SecretHandler,
	})
	return jwtMiddleware.Handler(next)
}
