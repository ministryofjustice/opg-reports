# setup base docker image building the application
FROM golang:1.24.4-alpine@sha256:68932fa6d4d4059845c8f40ad7e654e626f3ebd3706eef7846f319293ab5cb7a AS build
WORKDIR /app
COPY --link go.sum go.mod ./
COPY --link report ./report
RUN mkdir -p /gobin
ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin PATH=/gobin:$PATH
RUN update-ca-certificates
RUN apk --update --no-cache add build-base
RUN apk --update --no-cache add "sqlite=>=3.47"
RUN env CGO_ENABLED=1 go mod tidy
RUN env CGO_ENABLED=1 go build -o /gobin/api ./report/cmd/api/ && echo "OK"

# setup the development image
FROM build AS development
EXPOSE 8081
ENV CGO_ENABLED=1 GOOS=linux
RUN mkdir -p /gobin
WORKDIR /gobin
COPY --from=build /gobin/api /gobin
RUN rm -Rf /app
RUN apk --update --no-cache add "sqlite=>=3.47"
ENTRYPOINT ["./api"]

# setup production image
FROM alpine:3.22.0@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 AS production
EXPOSE 8081
ENV CGO_ENABLED=1 GOOS=linux
RUN mkdir -p /gobin
WORKDIR /gobin
COPY --from=build /gobin/api /gobin
RUN apk --update --no-cache add "sqlite=>=3.47"
ENTRYPOINT ["./api"]
