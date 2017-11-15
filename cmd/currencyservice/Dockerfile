FROM golang:1.9.2-stretch

RUN apt-get update && apt-get -y upgrade


# Run code
WORKDIR $GOPATH
CMD ["./bin/currencyservice"]
