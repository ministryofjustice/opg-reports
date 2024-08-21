# base requires gcc, sqlite
# install sqlc
FROM golang:1.22-alpine AS base
WORKDIR /app
ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin PATH=$GOBIN:$PATH
RUN mkdir -p /gobin
RUN update-ca-certificates
RUN apk --update --no-cache add build-base
RUN apk --update --no-cache add "sqlite=~3.45"
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
### create folders
RUN mkdir -p /go/bin/api/
RUN mkdir -p /app/builds/api/github_standards/data/
RUN mkdir -p /go/bin/api/github_standards/data/
RUN mkdir -p /app/builds/api/aws_costs/data/
RUN mkdir -p /go/bin/api/aws_costs/data/

# build stage - runs seeders
FROM base AS build
WORKDIR /app
# setup path vars - make sure gobin is included
ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin PATH=$GOBIN:$PATH
COPY . .
### update sqlc generators
WORKDIR /app/datastore/github_standards/
RUN sqlc generate
WORKDIR /app/datastore/aws_costs/
RUN sqlc generate
### copy schema
WORKDIR /app/datastore/github_standards/
RUN cp -p github_standards*.sql /go/bin/api/github_standards/
WORKDIR /app/datastore/aws_costs/
RUN cp -p aws_costs*.sql /go/bin/api/aws_costs/
### copy data
WORKDIR /app/builds/api/github_standards/data/
RUN cp -p *.json /go/bin/api/github_standards/data/
WORKDIR /app/builds/api/aws_costs/data/
RUN cp -p *.json /go/bin/api/aws_costs/data/
### build binaries
WORKDIR /app
RUN go build -o /go/bin/api/seed_cmd ./commands/seed/main.go
RUN go build -o /go/bin/api/api_server ./servers/api/main.go
### SEEDERS
# github standards
RUN /go/bin/api/seed_cmd \
    -table github_standards \
    -schema /go/bin/api/github_standards/github_standards.sql \
    -db /go/bin/api/github_standards.db \
    -data "/go/bin/api/github_standards/data/*.json"
# aws costs
RUN /go/bin/api/seed_cmd \
    -table aws_costs \
    -schema /go/bin/api/aws_costs/aws_costs.sql \
    -db /go/bin/api/aws_costs.db \
    -data "/go/bin/api/aws_costs/data/*.json"
# test the import worked by getting the count
RUN sqlite3 /go/bin/api/github_standards.db "select count(*) from github_standards;"
RUN sqlite3 /go/bin/api/aws_costs.db "select count(*) from aws_costs;"


FROM build AS dev
EXPOSE 8081
WORKDIR /go/bin/api
ENV CGO_ENABLED=1 GOOS=linux
RUN apk --update --no-cache add "sqlite=~3.45"
COPY --from=build /go/bin/api/api_server ./
COPY --from=build /go/bin/api/github_standards.db ./
### remove raw data files
RUN rm -f /go/bin/api/github_standards/data/*.json
RUN rm -f /go/bin/api/aws_costs/data/*.json
ENTRYPOINT ["./api_server"]

# PRODUCTION
FROM alpine:3.20.2
WORKDIR /go/bin/
ENV CGO_ENABLED=1 GOOS=linux
RUN apk --update --no-cache add "sqlite=~3.45"
COPY --from=build /go/bin/api/api_server ./
COPY --from=build /go/bin/api/github_standards.db ./
### remove raw data files
RUN rm -f /go/bin/api/github_standards/data/*.json
RUN rm -f /go/bin/api/aws_costs/data/*.json
ENTRYPOINT ["./api_server"]
