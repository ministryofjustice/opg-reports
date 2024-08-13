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
RUN apk --update --no-cache add "sqlite=~3.45"
ENV CGO_ENABLED=1 GOOS=linux
COPY --from=build /go/bin/api* ./
RUN pwd
RUN ls -lh ./
RUN ls -lh ./github_standards/
ENTRYPOINT ["./api_server"]



# FROM alpine:3
# WORKDIR /go/bin
# RUN apk --update --no-cache add \
#     ca-certificates \
#     tzdata
# # patch vulnerabilities
# RUN apk upgrade --no-cache busybox libcrypto3 libssl3
# # copy everying in the api folder
# COPY --from=build /app/builds/api* ./
# ENTRYPOINT ["./api_server"]


# FROM golang:1.22-alpine AS base
# WORKDIR /app
# # setup gobin to top level dir so its not copied over
# ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin
# RUN mkdir -p /gobin
# RUN update-ca-certificates
# # add
# #   gcc for cgo (used by the sqlite3 lib) via build-base
# #   sqlite for the sqlite3 lib
# #   git for goreleaser
# RUN apk --update --no-cache add \
#     build-base \
#     sqlite \
#     git
# # add go releaser thats used in build process
# RUN go install github.com/goreleaser/goreleaser/v2@v2.1.0
# # install sqlc
# RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0

# FROM base AS build
# ARG AWS_VAULT_PROFILE="operator"
# ARG AWS_BUCKET="bucket"
# WORKDIR /app
# # setup path vars - make sure gobin is included
# ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin PATH=$GOBIN:$PATH
# COPY . .
# RUN make vars AWS_VAULT_PROFILE="${AWS_VAULT_PROFILE}" AWS_BUCKET="${AWS_BUCKET}"
# # call the make file build step
# RUN make go-build AWS_VAULT_PROFILE="${AWS_VAULT_PROFILE}" AWS_BUCKET="${AWS_BUCKET}"

# FROM build AS dev
# EXPOSE 8081
# WORKDIR /go/bin
# # copy everying in the api folder
# COPY --from=build /app/builds/api* ./
# ENTRYPOINT ["./api_server"]

# FROM alpine:3
# WORKDIR /go/bin
# RUN apk --update --no-cache add \
#     ca-certificates \
#     tzdata
# # patch vulnerabilities
# RUN apk upgrade --no-cache busybox libcrypto3 libssl3
# # copy everying in the api folder
# COPY --from=build /app/builds/api* ./
# ENTRYPOINT ["./api_server"]
