FROM golang:1.22-alpine
WORKDIR /air-pollution-service

RUN apk add --no-cache gcc musl-dev

COPY . .
RUN go mod tidy
RUN go get ./...

RUN go build -ldflags '-w -s' -a -o ./bin/server ./cmd/server

CMD ["/air-pollution-service/bin/server"]
EXPOSE 3333