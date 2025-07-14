# setup base docker image building the application
FROM golang:1.24.4-alpine@sha256:68932fa6d4d4059845c8f40ad7e654e626f3ebd3706eef7846f319293ab5cb7a AS build
COPY --link go.sum go.mod /app/
COPY --link ./report/ /app/report/
COPY --link ./report/cmd/front/templates/ /gobin/app/templates/
COPY --link ./report/cmd/front/local-assets/ /gobin/app/local-assets/
RUN update-ca-certificates
RUN apk --update --no-cache add build-base
RUN mkdir -p /gobin/app
RUN addgroup -S user && adduser -S -g user user
WORKDIR /app
ENV CGO_ENABLED=1 GOOS=linux GOBIN=/gobin/app/ PATH=/gobin:$PATH
RUN env CGO_ENABLED=1 go mod tidy
RUN env CGO_ENABLED=1 go build -o /gobin/app/front ./report/cmd/front/
RUN chown -Rf user:user /gobin/
RUN chmod -Rf 0777 /gobin/

# setup the development image
FROM build AS development
EXPOSE 8080
COPY --from=build /gobin/ /gobin
WORKDIR /gobin/app
ENV CGO_ENABLED=1 GOOS=linux
RUN chown -Rf user:user /gobin/
RUN chmod -Rf 0777 /gobin/
ENTRYPOINT ["/gobin/app/front"]

# setup production image
FROM alpine:3.22.0@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 AS production
# FROM scratch AS production
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /gobin/ /gobin
WORKDIR /gobin/app
ENV CGO_ENABLED=1 GOOS=linux
USER user
EXPOSE 8080
ENTRYPOINT ["/gobin/app/front"]
