#####################################
## Build Container
#####################################
FROM golang:1.10.3 AS build-env
ADD . /go/src/github.com/pcman312/hackathon
#RUN apk add build-base gcc glibc abuild binutils binutils-doc gcc-doc && \
RUN cd /go/src/github.com/pcman312/hackathon/component_tests && \
#    GOOS=linux GOARCH=amd64 go test -race -tags component -c -o /component-tests
    go test -race -tags component -c -o /component-tests

#####################################
## Final Container
#####################################
FROM debian
COPY --from=build-env /component-tests /

ENTRYPOINT /component-tests -test.v --convey-story --host hackathon:9090
#ENTRYPOINT ./component-tests -test.v --convey-story --host hackathon:9090
#ENTRYPOINT ls -l