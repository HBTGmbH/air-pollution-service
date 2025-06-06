FROM golang:1.24.3-alpine3.21 AS builder
WORKDIR /app
RUN apk add -q --no-cache tzdata \
 && mkdir -p ./build/usr/share && cp -R /usr/share/zoneinfo ./build/usr/share/
ADD . .
RUN go build -ldflags="-w -s -extldflags=-static" -a -o ./build/server

FROM scratch
ENV TZ="Europe/Berlin"
COPY --from=builder /app/build /
ENTRYPOINT ["/server"]