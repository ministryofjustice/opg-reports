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
RUN mkdir -p /go/bin/api/github_standards/csv/
RUN ls -lRh /app/builds/api/github_standards/csv*
# sqlc generation
RUN cd /app/datastore/github_standards/ && sqlc generate && cd -
# data migration
RUN cp /app/datastore/github_standards/github_standards*.sql /go/bin/api/github_standards/
RUN cp /app/builds/api/github_standards/csv/* /go/bin/api/github_standards/csv/ || echo "skip csvs"
RUN go build -o /go/bin/api/api_server ./servers/api/main.go
RUN go build -o /go/bin/api/seed_cmd ./commands/seed/main.go
# debug info
RUN pwd
RUN ls -lRh /go/bin/
# run the db seeding
RUN /go/bin/api/seed_cmd \
    -table github_standards \
    -schema /go/bin/api/github_standards/github_standards.sql \
    -db /go/bin/api/github_standards.db \
    -csv "/go/bin/api/github_standards/csv/*.csv"
RUN sqlite3 /go/bin/api/github_standards.db "select count(*) from github_standards;"
RUN ls -lRh /go/bin/


FROM build AS dev
EXPOSE 8081
WORKDIR /go/bin/api
ENV CGO_ENABLED=1 GOOS=linux
RUN apk --update --no-cache add "sqlite=~3.45"
COPY --from=build /go/bin/api/api_server ./
COPY --from=build /go/bin/api/github_standards.db ./
RUN pwd
RUN ls -lRh ./
ENTRYPOINT ["./api_server"]

# PRODUCTION
FROM alpine:3.20.2
WORKDIR /go/bin
ENV CGO_ENABLED=1 GOOS=linux
RUN apk --update --no-cache add "sqlite=~3.45"
RUN mkdir -p ./github_standards
COPY --from=build /app/datastore/github_standards/github_standards*.sql ./github_standards/
COPY --from=build /app/builds/api/github_standards/csv* ./github_standards/
COPY --from=build /go/bin/api* ./
RUN pwd
RUN ls -lh ./
RUN ls -lh ./github_standards/
ENTRYPOINT ["./api_server"]
