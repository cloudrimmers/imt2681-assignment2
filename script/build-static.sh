#!/bin/sh

static_build() {
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -v -o $1 $2
}

static_build ./bin/rimbot ./cmd/rimbot
static_build ./bin/fixerworker ./cmd/fixerworker
static_build ./bin/currencyservice ./cmd/currencyservice
