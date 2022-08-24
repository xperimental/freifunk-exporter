# freifunk-exporter

This is a small program to extract some metrics from a Freifunk network. It gathers the data from the JSON file typically used to also populate the Meshviewer map. The data is then made available in a format readable by a [Prometheus](https://prometheus.io) server.

## Usage

To compile the tool, first check that you have a recent Go installation (1.19 at the time of this writing). Then clone the repository and run `make` which will run the tests and compile the binary.

The tool accepts a few parameters, only the `--source-url` is required:

```plain
Usage of freifunk-exporter:
      --addr string               Address to listen on. (default ":9295")
      --cache-interval duration   Interval for local caching of Meshviewer data. (default 3m0s)
      --source-url string         URL to Meshviewer JSON file.
```

So, for example to get the metrics for Freifunk Bodensee run:

```bash
freifunk-exporter --source-url https://meta.ffbsee.net/data/meshviewer.json
```
