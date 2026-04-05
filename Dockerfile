FROM golang:1.25.0

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o app ./cmd/api

CMD ["./app"]