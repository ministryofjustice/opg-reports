# setup base docker image building the application
FROM golang:1.26.10-alpine@sha256:a2012b4a1d1f9f7e9f44e0d39b805a901e856b1dd80b4e6ee4338db547050308 AS build
COPY --link go.sum go.mod /app/
COPY --link report /app/report
COPY --link builds/database /gobin/app/database
RUN update-ca-certificates
RUN apk --update --no-cache add build-base
RUN apk --update --no-cache add "sqlite=>=3.47"
RUN mkdir -p /gobin/app
RUN addgroup -S user && adduser -S -g user user
WORKDIR /app
ENV GOBIN=/gobin PATH=/gobin:$PATH
ENV CGO_ENABLED=1
RUN go mod download
RUN go build -ldflags="-w -s" -o /gobin/app/api ./report/cmd/api/
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
FROM alpine:3.23.7@sha256:ff7f0c727dd796aa36787bfcd9a64bf35a6ce72a73451b3359fde849b6703d63 AS production
EXPOSE 8081
# ssl certs
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# users & groups
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

# Patch Vulnerabilities
RUN apk upgrade --no-cache zlib

# app
COPY --from=build /gobin/ /gobin
WORKDIR /gobin/app
USER user

ENTRYPOINT ["/gobin/app/api"]
