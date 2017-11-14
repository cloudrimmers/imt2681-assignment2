FROM golang:1.9.2-stretch

RUN apt-get update && apt-get -y upgrade


# Install dependencies
RUN go get github.com/gorilla/mux \
    go get github.com/subosito/gotenv \
    go get gopkg.in/mgo.v2 \
    go get github.com/cloudrimmers/imt2681-assignment3

# Run code
WORKDIR $GOPATH
CMD ["./bin/rimbot"]
