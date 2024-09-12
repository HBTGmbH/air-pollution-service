FROM golang:1.22-alpine
WORKDIR /air-pollution-service

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/server ./cmd/server

CMD ["/air-pollution-service/bin/server"]
EXPOSE 3333