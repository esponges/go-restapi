package middleware

import (
	"fmt"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Logging middleware: %s %s ", r.Method, r.URL.Path)
		f(w, r)
	}
}

func Method(m string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			fmt.Printf("Method %s ", m)
			f(w, r)
		}
	}
}

func Wrapper(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		f = middleware(f)
	}
	return f
}
