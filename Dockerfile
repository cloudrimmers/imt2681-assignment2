FROM alpine:latest

MAINTAINER Jonas Solsvik <jonas.solsvik@gmail.com>

WORKDIR "/opt"

ADD .docker_build/imt2681-assignment2 /opt/bin/imt2681-assignment2

CMD ["/opt/bin/imt2681-assignment2"]