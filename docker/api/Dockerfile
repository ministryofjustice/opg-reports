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
#--fork-remove-start
# aws_costs
RUN mkdir -p /app/builds/api/aws_costs/data/
RUN mkdir -p /go/bin/api/aws_costs/data/
# aws_uptime
RUN mkdir -p /app/builds/api/aws_uptime/data/
RUN mkdir -p /go/bin/api/aws_uptime/data/
#--fork-remove-end

# build stage - runs seeders
FROM base AS build
WORKDIR /app
# setup path vars - make sure gobin is included
ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin PATH=$GOBIN:$PATH
COPY . .
###
# - update sqlc generators
# - copy schema
# - copy data
WORKDIR /app/datastore/github_standards/
RUN sqlc generate
WORKDIR /app/datastore/github_standards/
RUN cp -p github_standards*.sql /go/bin/api/github_standards/
WORKDIR /app/builds/api/github_standards/data/
RUN cp -p *.json /go/bin/api/github_standards/data/
#--fork-remove-start
# aws_costs
WORKDIR /app/datastore/aws_costs/
RUN sqlc generate
WORKDIR /app/datastore/aws_costs/
RUN cp -p aws_costs*.sql /go/bin/api/aws_costs/
WORKDIR /app/builds/api/aws_costs/data/
RUN cp -p *.json /go/bin/api/aws_costs/data/
# aws_uptime
WORKDIR /app/datastore/aws_uptime/
RUN sqlc generate
WORKDIR /app/datastore/aws_uptime/
RUN cp -p aws_uptime*.sql /go/bin/api/aws_uptime/
# WORKDIR /app/builds/api/aws_uptime/data/
# RUN cp -p *.json /go/bin/api/aws_uptime/data/
#--fork-remove-end
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
#--fork-remove-start
# aws_costs
RUN /go/bin/api/seed_cmd \
    -table aws_costs \
    -schema /go/bin/api/aws_costs/aws_costs.sql \
    -db /go/bin/api/aws_costs.db \
    -data "/go/bin/api/aws_costs/data/*.json"
RUN sqlite3 /go/bin/api/aws_costs.db "select count(*) from aws_costs;"
# aws_uptime
RUN /go/bin/api/seed_cmd \
    -table aws_uptime \
    -schema /go/bin/api/aws_uptime/aws_uptime.sql \
    -db /go/bin/api/aws_uptime.db \
    -data "/go/bin/api/aws_uptime/data/*.json"
RUN sqlite3 /go/bin/api/aws_uptime.db "select count(*) from aws_uptime;"
#--fork-remove-end
# test the import worked by getting the count
RUN sqlite3 /go/bin/api/github_standards.db "select count(*) from github_standards;"



FROM build AS dev
EXPOSE 8081
WORKDIR /go/bin/api
ENV CGO_ENABLED=1 GOOS=linux
RUN apk --update --no-cache add "sqlite=~3.45"
COPY --from=build /go/bin/api/api_server ./
COPY --from=build /go/bin/api/github_standards.db ./
RUN rm -f /go/bin/api/github_standards/data/*.json
#--fork-remove-start
# aws_costs
COPY --from=build /go/bin/api/aws_costs.db ./
RUN rm -f /go/bin/api/aws_costs/data/*.json
# aws_uptime
COPY --from=build /go/bin/api/aws_uptime.db ./
RUN rm -f /go/bin/api/aws_uptime/data/*.json
#--fork-remove-end
ENTRYPOINT ["./api_server"]

# PRODUCTION
FROM alpine:3.20.2
WORKDIR /go/bin/
ENV CGO_ENABLED=1 GOOS=linux
RUN apk --update --no-cache add "sqlite=~3.45"
COPY --from=build /go/bin/api/api_server ./
COPY --from=build /go/bin/api/github_standards.db ./
RUN rm -f /go/bin/api/github_standards/data/*.json
#--fork-remove-start
# aws_costs
COPY --from=build /go/bin/api/aws_costs.db ./
RUN rm -f /go/bin/api/aws_costs/data/*.json
# aws_uptime
COPY --from=build /go/bin/api/aws_uptime.db ./
RUN rm -f /go/bin/api/aws_uptime/data/*.json
#--fork-remove-end
ENTRYPOINT ["./api_server"]
