FROM golang:latest as build

WORKDIR /app

COPY app/go.mod ./
RUN go mod download

COPY app/*.go ./
COPY app/*.html ./

RUN go build -o /uk_realtime_bus_map

EXPOSE 8080

CMD ["/uk_realtime_bus_map"]