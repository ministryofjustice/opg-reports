# setup base docker image building the application
FROM golang:1.23.4-alpine AS build
ARG LDFLAGS
WORKDIR /app
COPY --link info ./info
COPY --link internal ./internal
COPY --link models ./models
COPY --link seed ./seed
COPY --link servers ./servers
COPY --link Makefile ./Makefile
COPY --link go.mod ./go.mod
COPY --link go.sum ./go.sum
ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin PATH=/gobin:$PATH
RUN mkdir -p /gobin
RUN update-ca-certificates
RUN apk --update --no-cache add build-base
RUN apk --update --no-cache add "sqlite=~3.47"
RUN date && echo "${LDFLAGS}" | sed "s#github.com/ministryofjustice/opg-reports/info.##g" | sed "s#-X#\\n#g" | sed "s#'##g"
RUN env CGO_ENABLED=1 go mod tidy
RUN env CGO_ENABLED=1 go build -ldflags="${LDFLAGS}" -o /gobin/api ./servers/api/main.go && echo "OK"


# setup the development image
FROM build AS development
EXPOSE 8081
ENV CGO_ENABLED=1 GOOS=linux
RUN mkdir -p /gobin
WORKDIR /gobin
COPY --from=build /gobin/api /gobin
RUN rm -Rf /app
RUN apk --update --no-cache add "sqlite=~3.47"
ENTRYPOINT ["./api"]


# setup production image
FROM alpine:3.21.0
ENV CGO_ENABLED=1 GOOS=linux
RUN mkdir -p /gobin
WORKDIR /gobin
COPY --from=build /gobin/api /gobin
RUN apk --update --no-cache add "sqlite=~3.47"
ENTRYPOINT ["./api"]
