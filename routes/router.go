package routes

import (
	"github.com/gorilla/mux"
	//"net/http"
	"net/http"
)

func NewRouter(loadBalancer *LoadBalancer) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	var routes = loadBalancer.CreateRoutes()
	for _, route := range routes {
		var handler http.Handler = route.HandlerFunc
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return router
}
