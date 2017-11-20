FROM golang:1.9.2-stretch as builder

#
# 1. BUILD STATGE
#
ENV PROJECT_ROOT=/go/src/github.com/cloudrimmers/imt2681-assignment3
COPY . $PROJECT_ROOT
WORKDIR $PROJECT_ROOT

RUN CGO_ENABLED=0 GOOS=linux \
		go build -a \
				 -ldflags '-extldflags "-static"' \
				 -v \
				 -o ./bin/fixerworker \
				 ./cmd/fixerworker
#
# 2. SECOND FROM SCRATCH STAGE 
#
FROM scratch

ENV PROJECT_ROOT=/go/src/github.com/cloudrimmers/imt2681-assignment3
WORKDIR /
COPY --from=builder $PROJECT_ROOT/bin/fixerworker .


# Run code
CMD ["/fixerworker", "docker"]
