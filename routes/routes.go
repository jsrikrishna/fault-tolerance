package routes

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Resources",
		"GET",
		"/resources",
		Resources,
	},
	Route{
		"Status",
		"POST",
		"/status",
		RequestStatusHandler,
	},
	Route{
		"Server",
		"POST",
		"/server",
		AddServer,
	},
}
