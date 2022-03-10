# syntax=docker/dockerfile:1

FROM golang:1.17-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN go build -o /bago

FROM gcr.io/distroless/base-debian11

ENV DATABASE_URL=postgres://root:root@localhost:5432/root
ENV JWT_SECRET=secret
ENV GIN_MODE=release

COPY --from=builder /bago /bago

EXPOSE 8080

ENTRYPOINT ["/bago"]