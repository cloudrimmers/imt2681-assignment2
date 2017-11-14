#!/bin/sh

govendor sync
docker build --tag currencyservice:cs \
		     --file ./cmd/currencyservice/Dockerfile \
	         .
