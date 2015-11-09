package authware

import (
	"encoding/json"
	"fmt"
	"github.com/InteractiveLecture/serviceclient"
	"github.com/gorilla/mux"
	"net/http"
)

func New(next http.Handler, objectClass string, permissions ...string) http.Handler {
	return createHandlerFunc(objectClass, permissions...)
}

func createHandlerFunc(objectClass string, permissions ...string) http.HandlerFunc {
	result := func(w http.ResponseWriter, r *http.Request) {
		instance := serviceclient.GetInstance("acl-service")
		vars := mux.Vars(r)
		sid := "" //TODO sid aus jwt-token auslesen.
		id, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := instance.Get(fmt.Sprintf("/%s/%s/permissions?sid=%s", objectClass, id, sid))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		decoder := json.NewDecoder(resp.Body)
		var resultPermissions = make([]string, 0)
		err = decoder.Decode(resultPermissions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !checkResult(permissions, resultPermissions) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	return result
}

func checkResult(expected []string, current []string) bool {
	//TODO intersection between result and expected
	return true
}
