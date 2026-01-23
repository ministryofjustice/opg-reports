ARG TARGETOS TARGETARCH
# setup base docker image building the application
FROM golang:1.25.5-alpine@sha256:ac09a5f469f307e5da71e766b0bd59c9c49ea460a528cc3e6686513d64a6f1fb AS build
COPY --link go.sum go.mod /app/
COPY --link ./report/ /app/report/
COPY --link ./report/cmd/front/templates/ /gobin/app/templates/
COPY --link ./report/cmd/front/local-assets/ /gobin/app/local-assets/
COPY --link builds/front/govuk /gobin/app/govuk
RUN update-ca-certificates
RUN apk --update --no-cache add build-base
RUN mkdir -p /gobin/app
RUN addgroup -S user && adduser -S -g user user
WORKDIR /app
ENV GOBIN=/gobin PATH=/gobin:$PATH
ENV CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH}
RUN go mod download
RUN go build -ldflags="-w -s" -o /gobin/app/front ./report/cmd/front/
RUN chown -Rf user:user /gobin/
RUN chmod -Rf 0755 /gobin/

# setup the development image
FROM build AS development
EXPOSE 8080
COPY --from=build /gobin/ /gobin
WORKDIR /gobin/app
ENV CGO_ENABLED=0 GOOS=linux
RUN chown -Rf user:user /gobin/
RUN chmod -Rf 0755 /gobin/
ENTRYPOINT ["/gobin/app/front"]

# setup production image
# FROM alpine:3.22.0@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 AS production
FROM scratch AS production
EXPOSE 8080
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /gobin/ /gobin
WORKDIR /gobin/app
USER user
ENTRYPOINT ["/gobin/app/front"]
