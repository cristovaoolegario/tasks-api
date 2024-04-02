FROM golang:1.21.3-alpine3.18 as builder
WORKDIR /work

# Download module in a separate layer to allow caching for the Docker build
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o microservice ./cmd/api/main.go

FROM alpine:3.18.6
WORKDIR /bin
COPY --from=builder /work/microservice /bin/microservice
EXPOSE 8080
CMD /bin/microservice
