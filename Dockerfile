FROM golang:1.14 as builder

COPY go.* /src/
WORKDIR /src
RUN go mod download

ADD . /src

ENV CGO_ENABLED=0
RUN go build -o server cmd/main.go

# ------

FROM alpine:latest

COPY --from=builder /src/server /usr/local/bin/server

EXPOSE 3000

ENTRYPOINT ["server"]
