package server

import "net/http"

type Route struct {
	Path        string
	Middlewares Middlewares
	Controller  func(res http.ResponseWriter, req *http.Request)
	Methods     []string
}

type Routes map[string]Route