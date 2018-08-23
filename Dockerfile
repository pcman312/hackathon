# build stage
FROM golang:1.10.3-alpine AS build-env
ADD . /go/src/github.com/pcman312/hackathon
RUN cd /go/src/github.com/pcman312/hackathon && \
    go build -o hackathon

FROM alpine
COPY --from=build-env \
  # Source files
  /go/src/github.com/pcman312/hackathon/hackathon \
  /go/src/github.com/pcman312/hackathon/seelog.xml \
  /go/src/github.com/pcman312/hackathon/env \
  # Destination
  /

EXPOSE 80
EXPOSE 9090

ENTRYPOINT /hackathon