package config

type Configuration struct {
	Name string `json:name`
	Protocol string `json:"protocol"`
	BindTo string `json:"bindto"`
	ClientIdleTimeout string `json:"clientIdleTimeout"`
	BackendIdleTimeout string `json:"backendIdleTimeout"`
	BackendConnectionTimeout string `json: "backendConnectionTimeout"`
	PingInterval int `json:"pingInterval"`
	HealthcheckInterval int `json:"healthcheckInterval"`
	Servers []Server `json:"servers"`
}

type Server struct {
	Address string `json:"address"`
	Name string `json:"name"`
	Status bool
}
