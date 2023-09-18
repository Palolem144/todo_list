FROM golang:1.18

WORKDIR /go/src/

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download -x

COPY . /go/src

RUN go build cmd/app/main.go

EXPOSE 8080

CMD ["./main"]