package info

import (
	"os"
	"testing"
)

func TestParseInfo(t *testing.T) {
	tt := []struct {
		desc      string
		filename  string
		wantNodes int
		wantLinks int
	}{
		{
			desc:      "simple",
			filename:  "testdata/simple.json",
			wantNodes: 2,
			wantLinks: 1,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			file, err := os.Open(tc.filename)
			if err != nil {
				t.Fatalf("error opening file: %s", err)
			}
			defer file.Close()

			info, err := parseInfo(file)
			if err != nil {
				t.Fatalf("error parsing info: %s", err)
			}

			if len(info.Nodes) != tc.wantNodes {
				t.Errorf("got %d nodes, want %d", len(info.Nodes), tc.wantNodes)
			}

			if len(info.Links) != tc.wantLinks {
				t.Errorf("got %d links, want %d", len(info.Links), tc.wantLinks)
			}
		})
	}
}
