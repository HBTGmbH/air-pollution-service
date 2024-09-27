ARG GOVERSION=1.23.1-alpine3.20
FROM golang:$GOVERSION AS builder
RUN apk add -q --no-cache tzdata
WORKDIR /air-pollution-service
ADD go.mod go.sum main.go ./
ADD config ./config
ADD internal ./internal
ADD docs ./docs
RUN go build -ldflags="-w -s -extldflags=-static" -a -o ./bin/server

FROM scratch
ENV TZ="Europe/Berlin"
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /air-pollution-service/bin/server /air-pollution-service/bin/server
ENTRYPOINT ["/air-pollution-service/bin/server"]

