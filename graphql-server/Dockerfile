ARG VERSION

FROM golang:1.19 as build

## Build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD . .

RUN go generate ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=$VERSION" -o diseases-risk-factors cmd/main.go

## Deploy
FROM alpine:latest
WORKDIR /app

COPY --from=build /app/diseases-risk-factors .

EXPOSE 8080

ENTRYPOINT ["./diseases-risk-factors", "--port", "8080"]