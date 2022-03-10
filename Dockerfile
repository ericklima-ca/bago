# syntax=docker/dockerfile:1

FROM golang:1.17-alpine as builder
WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN go build -o /bago

FROM golang:1.17

ENV DATABASE_URL=postgres://root:root@localhost:5432/root

COPY --from=builder /bago /bago

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/bago"]