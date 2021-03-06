package config

type Configuration struct {
	Name                     string `json:name`
	Protocol                 string `json:"protocol"`
	BindTo                   string `json:"bindto"`
	BindToStatusServer             string `json: "bindToStatuSserver"`
	ClientIdleTimeout        string `json:"clientIdleTimeout"`
	BackendIdleTimeout       string `json:"backendIdleTimeout"`
	BackendConnectionTimeout string `json: "backendConnectionTimeout"`
	PingInterval             int `json:"pingInterval"`
	HealthcheckInterval      int `json:"healthcheckInterval"`
	StatusCounter            int `json:"status_counter"`
	Servers                  []Server `json:"servers"`
	Algorithm                string `json:"algorithm"`
}

type Server struct {
	Address        string `json:"address"`
	Name           string `json:"name"`
	Status         bool
	CurrentCounter int
	Dead           bool
	Weight         int `json:"weight"`
	UsedResources 	float64
	FreeResources	uint64
}
