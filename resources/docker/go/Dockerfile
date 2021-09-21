# syntax=docker/dockerfile:1

# BUILD stage
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

WORKDIR /app/cmd
COPY cmd/ ./

WORKDIR /app/pkg
COPY pkg/ ./

WORKDIR /app
RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o ./mongo-svc-quickstart ./cmd


# DEPLOY stage
FROM gcr.io/distroless/base-debian10

WORKDIR deploy

COPY --from=build /app/mongo-svc-quickstart /deploy/
EXPOSE 9090
CMD ["/deploy/mongo-svc-quickstart"]
