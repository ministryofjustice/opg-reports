# base requires gcc, sqlite
# install sqlc
FROM golang:1.22-alpine AS base
WORKDIR /app
ENV CGO_ENABLED=1 GOOS=linux
RUN apk --update --no-cache add build-base
RUN update-ca-certificates

FROM base AS build
WORKDIR /app
ENV CGO_ENABLED=1 GOOS=linux
COPY . .
RUN mkdir -p /go/bin/front
RUN go build -o /go/bin/front/front_server ./servers/front/main.go
RUN pwd
RUN ls -lh /go/bin/
RUN ls -lh /go/bin/front


FROM build AS dev
EXPOSE 8080
WORKDIR /go/bin/front
ENV CGO_ENABLED=1 GOOS=linux
RUN mkdir -p ./static/
RUN mkdir -p ./templates/
COPY --from=build /app/servers/front/templates/* ./templates/
COPY --from=build /app/servers/front/static/* ./static/
COPY --from=build /app/servers/front/config.json ./
COPY --from=build /go/bin/front* ./
RUN pwd
RUN ls -lh ./
ENTRYPOINT ["./front_server"]

# PRODUCTION
FROM alpine:3.20.2
WORKDIR /go/bin
ENV CGO_ENABLED=1 GOOS=linux
RUN mkdir -p ./static/
RUN mkdir -p ./templates/
COPY --from=build /app/servers/front/templates/* ./templates/
COPY --from=build /app/servers/front/static/* ./static/
COPY --from=build /app/servers/front/config.json ./
COPY --from=build /go/bin/front* ./
RUN pwd
RUN ls -lh ./
ENTRYPOINT ["./front_server"]
