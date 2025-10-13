FROM golang:1.25.2-alpine3.22 AS builder
WORKDIR /app
RUN apk add -q --no-cache tzdata \
 && mkdir -p ./build/usr/share && cp -R /usr/share/zoneinfo ./build/usr/share/
ADD . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags '-w -s -extldflags "-static" -buildid=' -o ./build/server

FROM scratch
ENV TZ="Europe/Berlin"
COPY --from=builder /app/build /
ENTRYPOINT ["/server"]