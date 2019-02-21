package info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/FreifunkBremen/yanic/output/meshviewer"
)

// Nodes contains information about Freifunk network nodes.
type Nodes meshviewer.NodesV1

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
