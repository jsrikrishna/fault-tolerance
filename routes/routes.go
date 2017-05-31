package routes

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (loadBalancer *LoadBalancer) CreateRoutes() Routes {
	var routes = Routes{
		Route{
			"Resources",
			"GET",
			"/resources",
			loadBalancer.Resources,
		},
		Route{
			"Status",
			"POST",
			"/status",
			loadBalancer.RequestStatusHandler,
		},
		Route{
			"Server",
			"POST",
			"/server",
			loadBalancer.AddServer,
		},
	}
	return routes
}

