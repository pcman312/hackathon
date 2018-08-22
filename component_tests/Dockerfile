FROM golang:1.10.3

ADD . /go/src/github.com/pcman312/hackathon/component_tests

RUN cd /go/src/github.com/pcman312/hackathon/component_tests && go test -race -tags component -c -o /go/component-tests

ENTRYPOINT ["/go/component-tests", "-test.v", "--convey-story", "--host", "hackathon:9090"]