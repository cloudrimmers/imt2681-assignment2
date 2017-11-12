#!/bin/sh

go build -v -o ./bin/rimbot ./cmd/rimbot
go build -v -o ./bin/fixerworker ./cmd/fixerworker
go build -v -o ./bin/currencyservice ./cmd/currencyservice
