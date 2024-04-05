FROM golang:1.21.3-alpine3.18 as builder
WORKDIR /work
RUN apk update && apk add gcc librdkafka-dev zstd-libs libsasl lz4-dev libc-dev musl-dev

# Download module in a separate layer to allow caching for the Docker build
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd ./cmd
COPY internal ./internal

RUN go build -tags musl -o microservice ./cmd/api/rest/main.go

FROM alpine:3.18.6
WORKDIR /bin
RUN apk add --no-cache librdkafka
COPY --from=builder /work/microservice /bin/microservice
EXPOSE 8080
CMD /bin/microservice
