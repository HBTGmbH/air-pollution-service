FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app
RUN apk add -q --no-cache tzdata \
 && mkdir ./build \
 && mkdir -p ./build/usr/share && cp -R /usr/share/zoneinfo ./build/usr/share/
ADD go.mod go.sum main.go ./
ADD config ./config
ADD internal ./internal
ADD docs ./docs
RUN go build -ldflags="-w -s -extldflags=-static" -a -o ./build/server

FROM scratch
ENV TZ="Europe/Berlin"
COPY --from=builder /app/build /
ENTRYPOINT ["/server"]

