FROM golang:alpine

WORKDIR /app
COPY . /app

RUN go build -o /podarc main.go

ENTRYPOINT ["/podarc"]
