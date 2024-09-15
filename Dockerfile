FROM golang:1.22 as gobuild
ARG VERSION=latest

WORKDIR /air-pollution-service
ADD go.mod go.sum main.go ./
ADD config ./config
ADD internal ./internal
ADD docs ./docs

RUN go build -ldflags '-w -s' -a -o ./bin/server

FROM gcr.io/distroless/base

COPY --from=gobuild /air-pollution-service/bin/server /air-pollution-service/bin/server

ENTRYPOINT ["/air-pollution-service/bin/server"]
