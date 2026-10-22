package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/risicoding/nethawk/internal/protocol"
)

const PORT = 6363

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func GETROOT(ip string) (protocol.ServerWithLatency, error) {
	start := time.Now()
	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d", ip, PORT))
	duration := time.Since(start)

	if err != nil {
		return protocol.ServerWithLatency{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return protocol.ServerWithLatency{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var res protocol.Server

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return protocol.ServerWithLatency{}, err
	}

	if res.Name != "nethawk" || res.Status != "running" {
		return protocol.ServerWithLatency{}, fmt.Errorf("unreachable or not a nethawk server")
	}

	return protocol.ServerWithLatency{Server: protocol.Server{Name: res.Name, Config: res.Config, Status: res.Status, Ips: []string{ip}}, Latency: int(duration.Milliseconds())}, nil
}

func GETSCAN(ip string) ([]protocol.ServerWithLatency, error) {
	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d/scan", ip, PORT))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var res []protocol.ServerWithLatency

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
