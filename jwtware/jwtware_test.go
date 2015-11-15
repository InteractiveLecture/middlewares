package jwtware

import (
	"net/http"
	"testing"

	"github.com/InteractiveLecture/serviceclient/test"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestSecretGetter(t *testing.T) {
	controller := gomock.NewController(t)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("mysecret"))

	})
	server := serviceclienttest.Service(controller, "authentication-service", handler)
	defer server.Close()
	defer controller.Finish()
	result, err := SecretHandler(nil)
	assert.Nil(t, err)

	assert.Equal(t, "mysecret", string(result.([]byte)[:]))

}
