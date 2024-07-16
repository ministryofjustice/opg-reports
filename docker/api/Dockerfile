FROM golang:1.22-alpine AS base
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux
RUN update-ca-certificates

FROM base AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /go/bin/api services/api/main.go

FROM build as dev
EXPOSE 8081
WORKDIR /go/bin
COPY --from=build /go/bin/api api
COPY --from=build /app/services/api/data data
ENTRYPOINT ["./api"]

FROM alpine:3
WORKDIR /go/bin
RUN apk --update --no-cache add \
    ca-certificates \
    tzdata

# Patch vulnerabilities
RUN apk upgrade --no-cache busybox libcrypto3 libssl3

COPY --from=build /go/bin/api api
ENTRYPOINT ["./api"]
