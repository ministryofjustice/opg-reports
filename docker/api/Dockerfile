ARG TARGETOS TARGETARCH
# setup base docker image building the application
FROM golang:1.25.3-alpine@sha256:aee43c3ccbf24fdffb7295693b6e33b21e01baec1b2a55acc351fde345e9ec34 AS build
COPY --link go.sum go.mod /app/
COPY --link report /app/report
COPY --link builds/api/databases /gobin/app/databases
RUN update-ca-certificates
RUN apk --update --no-cache add build-base
RUN apk --update --no-cache add "sqlite=>=3.47"
RUN mkdir -p /gobin/app
RUN addgroup -S user && adduser -S -g user user
WORKDIR /app
ENV GOBIN=/gobin PATH=/gobin:$PATH
ENV CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH}
RUN go mod download
RUN go build -ldflags="-w -s" -o /gobin/app/api ./report/cmd/api/
# WORKDIR /gobin/app
# deal with linked binaries - sqlite package needs cgo.. grab the requirements and install them to ./lib/
# RUN ldd api | tr -s [:blank:] '\n' | grep ^/ | xargs -I % install -D % /gobin/app/%
RUN chown -Rf user:user /gobin/
RUN chmod -Rf 0755 /gobin/

# setup the development image
FROM build AS development
EXPOSE 8081
COPY --from=build /gobin/ /gobin
WORKDIR /gobin/app
RUN chown -Rf user:user /gobin/
RUN chmod -Rf 0755 /gobin/
ENTRYPOINT ["/gobin/app/api"]

# setup production image
FROM alpine:3.22.2@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412 AS production
# FROM scratch AS production
EXPOSE 8081
# ssl certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# users & groups
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
# app
COPY --from=build /gobin/ /gobin
# needs the linked files, can check using ldd (see the process above)
# COPY --from=build /gobin/app/lib/* /gobin/app/
WORKDIR /gobin/app
USER user
ENTRYPOINT ["/gobin/app/api"]
