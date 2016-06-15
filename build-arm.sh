#!/bin/bash -e

readonly TARGET=freifunk-exporter
readonly GO_PACKAGE=github.com/xperimental/freifunk-exporter

CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -a -tags netgo -ldflags "-w" -o "${TARGET}" "${GO_PACKAGE}"
