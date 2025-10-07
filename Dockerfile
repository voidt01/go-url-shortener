FROM golang:1.22.2

WORKDIR /usr/src/app

EXPOSE 4000

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o app ./cmd/api

CMD ["./app"]
