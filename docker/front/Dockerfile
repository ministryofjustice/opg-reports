FROM golang:1.22-alpine AS base
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux
RUN update-ca-certificates

FROM base AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /go/bin/front services/front/main.go

FROM build AS dev
EXPOSE 8080
WORKDIR /go/bin
COPY --from=build /go/bin/front front
COPY --from=build /app/services/front/config*.json ./
COPY --from=build /app/services/front/assets assets
COPY --from=build /app/services/front/templates templates
ENTRYPOINT ["./front"]

FROM alpine:3
WORKDIR /go/bin
RUN apk --update --no-cache add \
    ca-certificates \
    tzdata
# Patch vulnerabilities
RUN apk upgrade --no-cache busybox libcrypto3 libssl3
COPY --from=build /go/bin/front front
COPY --from=build /app/services/front/config*.json ./
COPY --from=build /app/services/front/assets assets
COPY --from=build /app/services/front/templates templates
ENTRYPOINT ["./front"]
