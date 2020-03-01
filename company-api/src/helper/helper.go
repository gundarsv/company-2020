package helper

import (
	"log"
	"net/http"
	"strings"
)

const (
	NilString = "NIL"
)

const (
	NoDatabaseError            = 0
	NotFoundDatabaseError      = 1
	RiskDatabaseError          = 2
	NoParametersDatabaseError  = 3
	AlreadyExistsDatabaseError = 4
)

func InitLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, r.Method)
		next.ServeHTTP(w, r)
	})
}

func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
		} else {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
			next.ServeHTTP(w, r)
		}
	})
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func GetStatusCode(errorCode int) int {
	switch errorCode {
	case 0:
		return http.StatusOK
	case 1:
		return http.StatusNotFound
	case 2:
		return http.StatusInternalServerError
	case 3:
		return http.StatusUnprocessableEntity
	case 4:
		return http.StatusForbidden
	}

	return 500
}
