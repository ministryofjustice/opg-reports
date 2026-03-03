# setup base docker image building the application
FROM golang:1.26.0-alpine@sha256:d4c4845f5d60c6a974c6000ce58ae079328d03ab7f721a0734277e69905473e5 AS build
COPY --link go.sum go.mod /app/
COPY --link report /app/report
COPY --link builds/templates /gobin/app/templates/
COPY --link builds/web /gobin/app/web/
COPY --link builds/govuk /gobin/app/govuk/
RUN update-ca-certificates
RUN apk --update --no-cache add build-base
RUN mkdir -p /gobin/app
RUN addgroup -S user && adduser -S -g user user
WORKDIR /app
ENV GOBIN=/gobin PATH=/gobin:$PATH
ENV CGO_ENABLED=0
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
FROM scratch AS production
EXPOSE 8080
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /gobin/ /gobin
WORKDIR /gobin/app
USER user
ENTRYPOINT ["/gobin/app/front"]
