# syntax=docker/dockerfile:1
# Dockerfile file for https://github.com/gostkov/simple_slow_http_server
FROM golang:1.19
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
#RUN CGO_ENABLED=0 GOOS=linux go build -o /simple_slow_http_server
RUN go build -o /simple_slow_http_server
EXPOSE 8080
CMD ["/simple_slow_http_server"]