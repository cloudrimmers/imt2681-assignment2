#!/bin/sh
go install ./cmd/fixerworker
go install ./cmd/webhookserver
heroku local