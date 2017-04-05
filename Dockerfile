FROM golang

MAINTAINER Justin Kiang <justin@bringhub.com>

ADD . /go/src/github.com/milodigital/cron-runner
WORKDIR /go/src/github.com/milodigital/cron-runner

CMD CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /tmp/cron-runner