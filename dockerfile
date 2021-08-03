# First stage
FROM golang:1.16-alpine as build

WORKDIR /goldcrest

COPY go.mod go.sum ./
COPY cmd/ ./cmd/
COPY protocol/ ./protocol/
COPY proxy/ ./proxy/

RUN go mod download
RUN mkdir -p bin
RUN go build -o bin/goldcrest cmd/server/server.go

# Second stage
FROM alpine:latest as run

COPY --from=build /goldcrest/bin/goldcrest /usr/local/bin/goldcrest

COPY resources/docker/config.yaml /etc/goldcrest/config.yaml
VOLUME /etc/goldcrest

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["/usr/local/bin/goldcrest"]
CMD ["-c", "/etc/goldcrest/config.yaml"]
