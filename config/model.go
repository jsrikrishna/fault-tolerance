package config

type Configuration struct {
	Name string `json:name`
	Protocol string `json:"protocol"`
	BindTo string `json:"bindto"`
	Servers []Server `json:"servers"`
}

type Server struct {
	Address string `json:"address"`
	Name string `json:"name"`
}
