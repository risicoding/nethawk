package protocol

type Server struct {
	Name   string       `json:"name"`
	Status string       `json:"status"`
	Config DeviceConfig `json:"config"`
	Ips    []string     `json:"ips"`
}

type ServerWithLatency struct {
	Server
	Latency int `json:"latency"`
}

type Device struct {
	DeviceConfig
	Ips []string `json:"ips"`
}

type DeviceConfig struct {
	Name  string `json:"name"`
	Uname string `json:"uname"`
	Alias string `json:"alias"`
}

type Interface struct {
	Name string `json:"name"`
	Ipv4 string `json:"ipv4"`
	Mask string `json:"mask"`
}

type ScanResult struct {
	Interface Interface `json:"interface"`
	Servers   []Server  `json:"servers"`
}
