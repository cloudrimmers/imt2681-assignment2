FROM golang:1.9.2-stretch

# RUN apt-get update && apt-get -y upgrade
ENV PROJECT_ROOT=$GOPATH/src/github.com/cloudrimmers/imt2681-assignment3
COPY . $PROJECT_ROOT
WORKDIR $PROJECT_ROOT

RUN go install -v ./cmd/fixerworker

CMD ["fixerworker", "docker"]
