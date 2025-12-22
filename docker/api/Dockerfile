ARG TARGETOS TARGETARCH
# setup base docker image building the application
FROM golang:1.25.5-alpine@sha256:ac09a5f469f307e5da71e766b0bd59c9c49ea460a528cc3e6686513d64a6f1fb AS build
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
FROM alpine:3.23.0@sha256:51183f2cfa6320055da30872f211093f9ff1d3cf06f39a0bdb212314c5dc7375 AS production
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
