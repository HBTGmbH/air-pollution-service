FROM golang:1.25.7-alpine3.23 AS builder
WORKDIR /app
RUN apk add -q --no-cache tzdata \
 && mkdir -p ./build/usr/share && cp -R /usr/share/zoneinfo ./build/usr/share/
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags '-w -s -extldflags "-static" -buildid=' -o ./build/server

FROM scratch
ENV TZ="Europe/Berlin"
COPY --from=builder /app/build /
ENTRYPOINT ["/server"]