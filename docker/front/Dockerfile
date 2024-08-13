FROM golang:1.22-alpine AS base
WORKDIR /app
# setup gobin to top level dir so its not copied over
ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin
RUN mkdir -p /gobin
RUN update-ca-certificates
# add
#   gcc for cgo (used by the sqlite3 lib) via build-base
#   sqlite for the sqlite3 lib
#   git for goreleaser
RUN apk --update --no-cache add \
    build-base \
    sqlite \
    git
# add go releaser thats used in build process
RUN go install github.com/goreleaser/goreleaser/v2@v2.1.0
# install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0

FROM base AS build
WORKDIR /app
# setup path vars - make sure gobin is included
ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin PATH=$GOBIN:$PATH
COPY . .
# call the make file build step
RUN make go-build

FROM build AS dev
EXPOSE 8081
WORKDIR /go/bin
# copy everying in the api folder
COPY --from=build /app/builds/front* ./
ENTRYPOINT ["./front_server"]

FROM alpine:3
WORKDIR /go/bin
RUN apk --update --no-cache add \
    ca-certificates \
    tzdata
# patch vulnerabilities
RUN apk upgrade --no-cache busybox libcrypto3 libssl3
# copy everying in the api folder
COPY --from=build /app/builds/front* ./
ENTRYPOINT ["./front_server"]
