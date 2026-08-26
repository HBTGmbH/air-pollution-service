FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.27.0-alpine@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS builder
RUN apk add --no-cache tzdata
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags '-w -s -extldflags "-static" -buildid=' -o server .

FROM scratch
ENV TZ="Europe/Berlin"
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/server /
USER 65534
ENTRYPOINT ["/server"]
