FROM golang:1.13.5-alpine3.10 AS builder
WORKDIR /build
#RUN adduser -u 10001 -D app-runner

ENV GOPROXY https://goproxy.cn
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o -a websocketSample .

FROM alpine:3.10 AS final
WORKDIR /app
COPY --from=builder /build/websocketSample /app/
#COPY --from=builder /build/config /app/config

#USER app-runner
EXPOSE 8090
ENTRYPOINT ["/app/websocketSample"]