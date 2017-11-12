#!/bin/sh

go install -v ./cmd/fixerworker
go install -v ./cmd/currencyservice
heroku local
