package info

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
    "meta": {
        "timestamp": "2016-06-18T22:00:04"
    },
    "nodes": [
        {
            "id": "de:ad:be:ef:40:38",
            "uptime": null,
            "flags": {
                "online": true,
                "gateway": false,
                "client": false
            },
            "name": "Router-1",
            "clientcount": 1,
            "hardware": null,
            "firmware": "ffbsee-0.0.6",
            "geo": [
                1.23,
                4.56
            ],
            "network": {
                "mac": "de:ad:be:ef:40:38"
            },
            "community": "bodensee"
        },
        {
            "id": "co:ff:ee:ba:be:83",
            "uptime": null,
            "flags": {
                "online": true,
                "gateway": true,
                "client": false
            },
            "name": "vpn1",
            "clientcount": 0,
            "hardware": null,
            "firmware": "server",
            "geo": null,
            "network": {
                "mac": "co:ff:ee:ba:be:83"
            },
            "community": "bodensee"
        }
    ],
    "links": [
        {
            "target": 1,
            "source": 0,
            "quality": "1, 1",
            "id": "deadbeef4038-coffeebabe83",
            "type": "vpn"
        }
    ]
}
`))
}

func TestGenerator(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(nodesHandler))

	tick := make(chan struct{})
	g := Generator(server.URL, tick)

	go func() {
		tick <- struct{}{}
		close(tick)
	}()

	nodes := <-g

	if nodes == nil {
		t.Errorf("got nil nodes")
	}

	if len(nodes.Nodes) != 2 {
		t.Errorf("got %d nodes, expected 2", len(nodes.Nodes))
	}

	if len(nodes.Links) != 1 {
		t.Errorf("got %d links, expected 1", len(nodes.Links))
	}
}
