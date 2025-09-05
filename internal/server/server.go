package server

import (
	"embed"

	"fmt"
	"log"
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api"
	"github.com/ycatan-w/himawari-api/internal/api/middleware"
	"github.com/ycatan-w/himawari-api/internal/output"
	"github.com/ycatan-w/himawari-api/internal/output/colors"
)

//go:embed "web/*"
var webFS embed.FS

type Server struct {
	Port int
}

type RouteDefinition struct {
	RouteName   string
	Methods     []string
	HandlerFunc http.HandlerFunc
	Handler     http.Handler
}

var routes = []RouteDefinition{
	{RouteName: "/api/login", Methods: []string{http.MethodPost}, HandlerFunc: api.LoginHandler},
	{RouteName: "/api/register", Methods: []string{http.MethodPost}, HandlerFunc: api.RegisterHandler},
	{RouteName: "/api/logout", Methods: []string{http.MethodPost}, HandlerFunc: api.LogoutHandler},
	{RouteName: "/api/events", Methods: []string{http.MethodGet, http.MethodPost}, HandlerFunc: middleware.AuthMiddleware(api.EventsHandler)},
	{RouteName: "/api/events/{id}", Methods: []string{http.MethodPut, http.MethodDelete}, HandlerFunc: middleware.AuthMiddleware(api.EventsByIdHandler)},
	{RouteName: "/api/logs", Methods: []string{http.MethodGet, http.MethodPost}, HandlerFunc: middleware.AuthMiddleware(api.LogsHandler)},
	{RouteName: "/api/logs/{id}", Methods: []string{http.MethodPut, http.MethodDelete}, HandlerFunc: middleware.AuthMiddleware(api.LogsByIdHandler)},
	{RouteName: "/web/", Methods: []string{http.MethodGet}, Handler: http.FileServer(http.FS(webFS))},
}

func New() *Server {
	return &Server{
		Port: 9740,
	}
}
func (s *Server) Run() {
	output.PrintHeader("Start himawari-server Server")
	mux := http.NewServeMux()

	output.PrintSubHeader("Register routes")
	for _, route := range routes {
		registerRoute(mux, route)
	}
	handler := middleware.CORS(mux)
	handler = logRequest(handler)

	log.Printf("[%s] Server will be running at %s\n", output.AppNameGreen(), colors.YellowUnderline(fmt.Sprintf("http://localhost:%d", s.Port)))
	output.PrintSubHeader("Listen to routes")
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), handler); err != nil {
		log.Fatal(err)
	}
}

func registerRoute(mux *http.ServeMux, route RouteDefinition) {
	allowed := make(map[string]bool, len(route.Methods))
	for _, m := range route.Methods {
		allowed[m] = true
	}
	mux.HandleFunc(route.RouteName, func(w http.ResponseWriter, r *http.Request) {
		if !middleware.ValidateRouteMethod(w, r, allowed, route.Methods) {
			return
		}
		if route.HandlerFunc != nil {
			route.HandlerFunc(w, r)
		} else if route.Handler != nil {
			route.Handler.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	for _, method := range route.Methods {
		log.Printf("[%s][%s] Register Route: %s", output.AppNameGreen(), colors.Yellow("route"), colors.BrightYellow(fmt.Sprintf("%-7s %s", method, route.RouteName)))
	}
}
