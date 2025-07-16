ARG TARGETOS TARGETARCH
# setup base docker image building the application
FROM golang:1.24.4-alpine@sha256:68932fa6d4d4059845c8f40ad7e654e626f3ebd3706eef7846f319293ab5cb7a AS build
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
RUN go mod download
RUN env CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /gobin/app/api ./report/cmd/api/
WORKDIR /gobin/app
# deal with linked binaries - sqlite package needs cgo..
RUN ldd api | tr -s [:blank:] '\n' | grep ^/ | xargs -I % install -D % /gobin/app/%
RUN chown -Rf user:user /gobin/
RUN chmod -Rf 0755 /gobin/

# setup the development image
FROM build AS development
EXPOSE 8081
COPY --from=build /gobin/ /gobin
WORKDIR /gobin/app
ENV CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH}
RUN chown -Rf user:user /gobin/
RUN chmod -Rf 0755 /gobin/
ENTRYPOINT ["/gobin/app/api"]

# setup production image
# FROM alpine:3.22.0@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 AS production
FROM scratch AS production
EXPOSE 8081
# ssl certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# users & groups
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
# app
COPY --from=build /gobin/ /gobin
# needs the linked files, can check using ldd
COPY --from=build /gobin/app/lib /
ENV CGO_ENABLED=1 CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH}
USER user
ENTRYPOINT ["/gobin/app/api"]
