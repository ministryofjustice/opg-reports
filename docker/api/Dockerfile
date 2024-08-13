# base requires gcc, sqlite
# install sqlc
FROM golang:1.22-alpine AS base
WORKDIR /app
ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin
RUN mkdir -p /gobin
RUN update-ca-certificates
RUN apk --update --no-cache add build-base
RUN apk --update --no-cache add "sqlite=~3.45"
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0


# sqlc generation
FROM base AS build
WORKDIR /app
# setup path vars - make sure gobin is included
ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin PATH=$GOBIN:$PATH
COPY . .
RUN mkdir -p /go/bin/api
# sqlc generation
RUN cd ./datastore/github_standards/ && sqlc generate && cd -
RUN go build -o /go/bin/api/api_server ./servers/api/main.go
RUN pwd
RUN ls -lh /go/bin/
RUN ls -lh /go/bin/api


FROM build AS dev
EXPOSE 8081
WORKDIR /go/bin/api
ENV CGO_ENABLED=1 GOOS=linux
RUN apk --update --no-cache add "sqlite=~3.45"
RUN mkdir -p ./github_standards
# copy schemas
COPY --from=build /app/datastore/github_standards/github_standards*.sql ./github_standards/
# copy bucket data
COPY --from=build /app/builds/api/github_standards/*.csv ./github_standards/
# copy binary
COPY --from=build /go/bin/api* ./
RUN pwd
RUN ls -lh ./
RUN ls -lh ./github_standards/
ENTRYPOINT ["./api_server"]

# PRODUCTION
FROM alpine:3.20.2
WORKDIR /go/bin
ENV CGO_ENABLED=1 GOOS=linux
RUN apk --update --no-cache add "sqlite=~3.45"
RUN mkdir -p ./github_standards
COPY --from=build /app/datastore/github_standards/github_standards*.sql ./github_standards/
COPY --from=build /app/builds/api/github_standards/*.csv ./github_standards/
COPY --from=build /go/bin/api* ./
RUN pwd
RUN ls -lh ./
RUN ls -lh ./github_standards/
ENTRYPOINT ["./api_server"]
