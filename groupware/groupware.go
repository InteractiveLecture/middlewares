package groupware

import (
	"net/http"
)

func New(next http.Handler, groups ...string) http.Handler {
	return http.HandlerFunc(checkHandlerFunc(next, groups...))
}

func checkHandlerFunc(next http.Handler, groups ...string) http.HandlerFunc {
	result := func(w http.ResponseWriter, r *http.Request) {
		actualGroups := make([]string, 0) //TODO read groups from context (parsed by jwtware)
		if !inGroup(actualGroups, groups) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
	return result
}

func inGroup(actualGroup []string, groups []string) bool {
	//TODO check if groups contain group
	return true
}
