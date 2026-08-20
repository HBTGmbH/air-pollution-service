FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.26.6-alpine@sha256:3889b425f035be855a72fb4755265311293b6d414521f0a519d819df32222d83 AS builder
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
