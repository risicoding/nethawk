package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/risicoding/nethawk/internal/protocol"
)

func main() {

}

const PORT = 5353

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func Query(ip string) (protocol.Server, error) {

	resp, err := httpClient.Get(fmt.Sprintf("http://%s:%d", ip, PORT))
	if err != nil {
		return protocol.Server{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return protocol.Server{}, fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var res protocol.Server

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return protocol.Server{}, err
	}

	if res.Name != "nethawk" || res.Status != "ok" {
		return protocol.Server{}, fmt.Errorf("not a nethawk server")
	}

	return res, nil
}
