FROM golang:1.7

RUN go get github.com/influxdata/influxdb/client/v2 && \
    go get github.com/shirou/gopsutil

ADD ./app.go ./app.go

CMD ["go", "run", "app.go"]
