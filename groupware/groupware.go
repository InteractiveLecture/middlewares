package groupware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"net/http"
)

type Options struct {
	Next      http.Handler
	Extractor ParameterExtractor
	Groups    map[string]bool
}

func New(options Options) http.Handler {
	return http.HandlerFunc(checkHandlerFunc(options))
}

func DefaultOptions(next http.Handler, groups ...string) Options {
	var groupMap = make(map[string]bool)
	for _, v := range groups {
		groupMap[v] = true
	}
	options := Options{
		next,
		defaultExtractor,
		groupMap,
	}
	return options
}

type ParameterExtractor func(r *http.Request) (groups []string)

func defaultExtractor(r *http.Request) (groups []string) {
	return context.Get(r, "user").(*jwt.Token).Claims["authorities"].([]string)
}

func checkHandlerFunc(options Options) http.HandlerFunc {
	result := func(w http.ResponseWriter, r *http.Request) {
		actualGroups := options.Extractor(r)
		if !inGroup(actualGroups, options.Groups) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if options.Next != nil {
			options.Next.ServeHTTP(w, r)
		}
	}
	return result
}

func inGroup(actualGroups []string, groups map[string]bool) bool {
	for _, v := range actualGroups {
		_, ok := groups[v]
		if ok {
			return true
		}
	}
	return false
}
