package middleware

import (
	"fmt"
	"net/http"
)

func Logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Logging middleware: %s %s ", r.Method, r.URL.Path)
		f(w, r)
	}
}
