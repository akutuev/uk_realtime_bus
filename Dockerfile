FROM golang:latest

WORKDIR /app

COPY app/go.mod ./
RUN go mod download

COPY . .

RUN cd app && go build -o /uk_realtime_bus

EXPOSE 8080

CMD ["/uk_realtime_bus"]