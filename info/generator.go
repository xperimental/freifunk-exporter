package info

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Meshinfo contains the information from the Meshviewer JSON file.
type Meshinfo struct {
	Links []Link `json:"links"`
	Nodes []Node `json:"nodes"`
	Meta  Meta   `json:"meta"`
}

// Node contains information about a single network node.
type Node struct {
	ID            string      `json:"node_id"`
	IsGateway     bool        `json:"is_gateway"`
	Uptime        string      `json:"uptime"`
	FirstSeen     string      `json:"firstseen"`
	LastSeen      string      `json:"lastseen"`
	Clients       int         `json:"clients"`
	ClientsWifi24 int         `json:"clients_wifi24"`
	ClientsWifi5  int         `json:"clients_wifi5"`
	ClientsOther  int         `json:"clients_other"`
	LoadAvg       float64     `json:"loadavg"`
	MemoryUsage   float64     `json:"memory_usage"`
	RootfsUsage   float64     `json:"rootfs_usage"`
	Firmware      Firmware    `json:"firmware"`
	Addresses     []string    `json:"addresses"`
	Contact       string      `json:"contact"`
	Autoupdater   Autoupdater `json:"autoupdater"`
	MAC           string      `json:"mac"`
	Hostname      string      `json:"hostname"`
	SiteCode      string      `json:"site_code"`
	VPN           bool        `json:"vpn"`
	Gateway       string      `json:"gateway"`
	Online        bool        `json:"is_online"`
	Location      Location    `json:"location"`
	Model         string      `json:"model"`
	// "gateway_nexthop": "-",
	// "nproc": 1,
}

// Link contains information about a link between two network nodes.
type Link struct {
	SourceID      string   `json:"source"`
	TargetID      string   `json:"target"`
	Type          LinkType `json:"type"`
	SourceAddress string   `json:"source_addr"`
	TargetAddress string   `json:"target_addr"`
	SourceQuality float64  `json:"source_tq"`
	TargetQuality float64  `json:"target_tq"`
}

// LinkType describes the type of a node link.
type LinkType string

const (
	// LinkTypeVPN signals a link between node and VPN gateway.
	LinkTypeVPN LinkType = "vpn"
	// LinkTypeOther signals a link between two nodes.
	LinkTypeOther LinkType = "other"
)

// Meta contains meta information about the Meshviewer data.
type Meta struct {
	Timestamp string `json:"timestamp"`
}

// Firmware contains information about the firmware used on a node.
type Firmware struct {
	Release string `json:"release"`
	Base    string `json:"base"`
}

// Autoupdater contains information about the state of the autoupdater on a node.
type Autoupdater struct {
	Enabled bool   `json:"enabled"`
	Branch  string `json:"branch"`
}

// Location contains the geographic location of a node.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetInfo will read node information from a HTTP URL.
func GetInfo(url string) (*Meshinfo, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to server: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-ok HTTP status: %d", res.StatusCode)
	}

	return parseInfo(res.Body)
}

func parseInfo(r io.Reader) (*Meshinfo, error) {
	var result Meshinfo
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		return nil, fmt.Errorf("parse error: %s", err)
	}

	return &result, nil
}
