#!/bin/sh
. ./.env
go install -v ./cmd/rimbot
go install -v ./cmd/fixerworker
go install -v ./cmd/currencyservice