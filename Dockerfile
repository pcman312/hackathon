# build stage
FROM golang:1.10.3 AS build-env
ADD . /go/src/github.com/pcman312/hackathon
RUN cd /go/src/github.com/pcman312/hackathon && go build -o hackathon

FROM debian

COPY --from=build-env /go/src/github.com/pcman312/hackathon/seelog.xml /seelog.xml
COPY --from=build-env /go/src/github.com/pcman312/hackathon/env /env
COPY --from=build-env /go/src/github.com/pcman312/hackathon/hackathon /hackathon

EXPOSE 9090

ENTRYPOINT "/hackathon"
