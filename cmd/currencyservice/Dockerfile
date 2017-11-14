FROM golang:1.9.2-stretch
LABEL maintainer="jonasjso@stud.ntnu.com"

RUN apt-get update && apt-get -y upgrade

# Copy and install code
COPY ../../. /code
WORKDIR /code
RUN ./script/install.sh

# Run code
WORKDIR $GOPATH
CMD ["./bin/rimbot"]