package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marllef/awesome-queue/pkg/utils/logger"
)

type Server interface {
	SetRoutes(routes Routes)
	SetPrefix(prefix string)
	GetPrefix() string
	SetLogger(logger *logger.Logger)
	GetLogger() *logger.Logger
	SetPort(port string)
	AddRoute(key string, route Route)
	GetRouter() *mux.Router
	Serve()
}

type server struct {
	routes Routes
	router *mux.Router
	port   string
	logger *logger.Logger
	prefix string
	Server
}

// Create a new server.
func NewServer() *server {
	return &server{
		routes: make(Routes),
		port:   "3005",
		prefix: "",
		router: mux.NewRouter(),
		logger: logger.Default(),
	}
}

// Get server routes.
func (s *server) GetRoutes() Routes {
	return s.routes
}

// Set server routes.
func (s *server) SetRoutes(routes Routes) {
	s.routes = routes
}

// Get server router
func (s *server) GetRouter() *mux.Router {
	return s.router
}

// Set server router
func (s *server) SetRouter(router *mux.Router)  {
	s.router = router
}

// Get server logger
func (s *server) GetLogger() *logger.Logger {
	return s.logger
}

// Set server logger.
func (s *server) SetLogger(logger *logger.Logger) {
	s.logger = logger
}

// Add a route in server.
func (s *server) AddRoute(key string, route Route) {
	s.routes[key] = route
}

// Set route prefix.
func (s *server) SetPrefix(prefix string) {
	s.prefix = prefix
}

// Get server prefix.
func (s *server) GetPrefix() string {
	return s.prefix
}

// Get server port.
func (s *server) GetPort(port string) string {
	return s.port
}

// Set server port.
func (s *server) SetPort(port string) {
	s.port = port
}

func (s *server) Serve() error {
	addr := fmt.Sprintf(":%s", s.port)
	s.logger.Infof("Servidor iniciado na porta 0.0.0.0:%s", s.port)

	for key, route := range s.routes {
		subRouter := s.router.Name(key).Subrouter()
		subRouter.Use(route.Middlewares...)

		path := fmt.Sprintf("%s%s", s.prefix, route.Path)

		subRouter.HandleFunc(path, route.Controller).Methods(route.Methods...)

		s.logger.Infof("New Route Available [%s]", path)
	}

	return http.ListenAndServe(addr, s.router)
}
