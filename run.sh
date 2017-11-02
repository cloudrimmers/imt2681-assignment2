#!/bin/sh
go install -v ./cmd/fixerworker
go install -v ./cmd/webhookserver
heroku local