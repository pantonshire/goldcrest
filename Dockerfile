FROM golang:1.16-alpine as build
WORKDIR /goldcrest
COPY go.mod go.sum ./
COPY cmd/ ./cmd/
COPY protocol/ ./protocol/
COPY proxy/ ./proxy/
RUN go mod download
RUN mkdir -p bin
RUN go build -o bin/goldcrest cmd/server/server.go

FROM alpine:latest as runtime
COPY --from=build /goldcrest/bin/goldcrest /usr/local/bin/goldcrest
RUN chmod 0755 /usr/local/bin/goldcrest
COPY resources/docker/config.yaml /etc/goldcrest/config.yaml
RUN chmod 0644 /etc/goldcrest/config.yaml
RUN addgroup -S goldcrest
RUN adduser -SDH -G goldcrest goldcrest
EXPOSE 8080
USER goldcrest:goldcrest
ENTRYPOINT ["/usr/local/bin/goldcrest", "-c", "/etc/goldcrest/config.yaml"]
