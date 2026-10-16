package server

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/risicoding/nethawk/internal/ip"
	"github.com/risicoding/nethawk/internal/protocol"
)

type Server protocol.Server
type Interface protocol.Interface

type ScanResult protocol.ScanResult

func Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		infs, err := ip.GetInterfaces()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var ips []string
		for _, inf := range infs {
			if inf.Name == "lo" || inf.Name == "docker0" {
				continue
			}

			addrs, err := inf.Addrs()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

			}

			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok {

					if ipv4 := ipNet.IP.To4(); ipv4 != nil {
						ips = append(ips, ipv4.String())
					}
				}
			}
		}

		deviceName := "Pentium Dual Core e6700; Fedora 6.19"

		var deviceUname string
		for _, s := range deviceName {
			deviceUname += strings.ToLower(string(s))
		}

		deviceUname = strings.Join(strings.Split(deviceUname, " "), "_")

		err = json.NewEncoder(w).Encode(Server{
			Name:        "nethawk",
			Ips:         ips,
			Status:      "ok",
			DeviceName:  deviceName,
			DeviceUname: deviceUname,
			DeviceAlias: "pc",
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/interface", func(w http.ResponseWriter, r *http.Request) {
		v, err := ip.GetInterfaces()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var infDatas []Interface
		for _, inf := range v {
			if inf.Name == "lo" || inf.Name == "docker" {
				continue
			}

			addrs, err := inf.Addrs()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok {

					if ipv4 := ipNet.IP.To4(); ipv4 != nil {
						data, err := ip.GetInterfaceData(inf.Name)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
						d := data
						infDatas = append(infDatas, Interface(d))
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(infDatas)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		value, err := ip.ScanAllInterfacesCached()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		var data []ScanResult

		for _, inf := range value {
			scResult := ScanResult{Interface: protocol.Interface(inf.Interface), Servers: inf.Servers}
			data = append(data, scResult)
		}

		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	err := http.ListenAndServe(":5353", nil)
	if err != nil {
		panic(err)
	}
}
