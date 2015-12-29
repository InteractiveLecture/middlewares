package jwtware

import (
	"net/http"
	"testing"

	"github.com/InteractiveLecture/serviceclient/test"
	"github.com/stretchr/testify/assert"
)

func TestSecretGetter(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("mysecret"))
	})
	server, _ := servicetest.Service("authentication-service", handler)
	defer server.Close()
	result, err := SecretHandler(nil)
	assert.Nil(t, err)
	assert.Equal(t, "mysecret", string(result.([]byte)[:]))
}
