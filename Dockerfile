# syntax=docker/dockerfile:1.3
FROM golang:1.21-bullseye AS BUILD

WORKDIR /src/autoproofcli

COPY go.mod go.sum ./
RUN go mod download &&  \
    go mod verify

COPY ./ ./

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -v -o /usr/local/bin/autoproofcli

FROM alpine:3.18
RUN addgroup --system autoproof && \
    adduser --system autoproof --ingroup autoproof

COPY --chown=autoproof:autoproof --from=BUILD /usr/local/bin/autoproofcli /usr/local/bin/autoproofcli
USER autoproof

ENTRYPOINT [ "/usr/local/bin/autoproofcli" ]