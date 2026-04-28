FROM golang:1.23-alpine AS build
WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /edge-agent ./cmd/agent

FROM alpine:3.20
WORKDIR /app
COPY --from=build /edge-agent /usr/local/bin/edge-agent

EXPOSE 8081
ENTRYPOINT ["/usr/local/bin/edge-agent"]
