FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o myapp .

FROM alpine:3.18

WORKDIR /

COPY --from=builder /app/myapp .

ENV TZ=Etc/UTC

ENTRYPOINT ["./myapp"]
