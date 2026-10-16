package protocol

type Server struct {
	Name        string   `json:"name"`
	Ips         []string `json:"ips"`
	Status      string   `json:"status"`
	DeviceName  string   `json:"devicename"`
	DeviceUname string   `json:"deviceuname"`
	DeviceAlias string   `json:"devicealias"`
}

type Interface struct {
	Name string `json:"name"`
	Ipv4 string `json:"ipv4"`
	Mask string `json:"mask"`
}

type ScanResult struct {
	Interface Interface `json:"interface"`
	Servers   []Server  `json"servers"`
}
