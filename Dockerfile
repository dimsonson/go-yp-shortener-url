FROM golang:latest

WORKDIR /app/shortener

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o cmd/shortener ./...

#CMD ["go", 'run' "cmd/shortener/main.go"]

RUN chmod +x cmd/shortener

CMD ["./cmd/shortener"]

