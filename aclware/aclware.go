package authware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/InteractiveLecture/serviceclient"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type PermissionFetcher func(id string, sid string) (map[string]interface{}, error)

type ParameterExtractor func(r *http.Request) (id string, sid string)

type Options struct {
	Permissions []string
	Fetcher     PermissionFetcher
	Next        http.Handler
	Extractor   ParameterExtractor
}

func defaultFetcher(id string, sid string) (result map[string]interface{}, err error) {
	instance := serviceclient.GetInstance("acl-service")
	var resp *http.Response
	resp, err = instance.Get(fmt.Sprintf("/objects/%s/permissions/%s", id, sid))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	return
}

func DefaultOptions(next http.Handler, idField string, permissions ...string) Options {
	fetcher := defaultFetcher
	extractor := func(r *http.Request) (id string, sid string) {
		vars := mux.Vars(r)
		user := context.Get(r, "user")
		sid = user.(*jwt.Token).Claims["user_name"].(string)
		id = vars[idField]
		return
	}
	return Options{
		Permissions: permissions,
		Next:        next,
		Fetcher:     fetcher,
		Extractor:   extractor,
	}
}

func New(options Options) http.Handler {
	return http.HandlerFunc(CreateHandlerFunc(options))
}

func CreateHandlerFunc(options Options) http.HandlerFunc {
	result := func(w http.ResponseWriter, r *http.Request) {
		id, sid := options.Extractor(r)
		resultPermissions, err := options.Fetcher(id, sid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !checkResult(options.Permissions, resultPermissions) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if options.Next != nil {
			options.Next.ServeHTTP(w, r)
		}
	}
	return result
}

func checkResult(expected []string, current map[string]interface{}) bool {
	for _, permission := range expected {
		if val, ok := current[permission]; !ok || !val.(bool) {
			return false
		}
	}
	return true
}
