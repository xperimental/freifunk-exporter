package info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Meta contains metadata about the node list.
type Meta struct {
	Timestamp string `json:"timestamp"`
}

// Nodes contains information about Freifunk network nodes.
type Nodes struct {
	Meta  Meta   `json:"meta"`
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

// Node contains information about one Freifunk network node.
type Node struct {
	ID        string      `json:"id"`
	Uptime    float64     `json:"uptime"`
	Flags     NodeFlags   `json:"flags"`
	Name      string      `json:"name"`
	Clients   int         `json:"clientcount"`
	Hardware  string      `json:"hardware"`
	Firmware  string      `json:"firmware"`
	Geo       []float64   `json:"geo"`
	Network   NodeNetwork `json:"network"`
	Community string      `json:"community"`
}

// NodeFlags contains flags communicating the purpose of the network node.
type NodeFlags struct {
	Gateway bool `json:"gateway"`
	Online  bool `json:"online"`
	Client  bool `json:"client"`
}

// NodeNetwork contains information about the network address of the node.
type NodeNetwork struct {
	MAC string `json:"mac"`
}

// Link contains information about a link between two network nodes.
type Link struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Source  int    `json:"source"`
	Target  int    `json:"target"`
	Quality string `json:"quality"`
}

// GetNodes will read node information from a HTTP URL.
func GetNodes(url string) (*Nodes, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-ok HTTP status: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can not read response: %s", err)
	}

	var result Nodes
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("can not read JSON: %s", err)
	}

	return &result, nil
}
