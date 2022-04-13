FROM golang

WORKDIR /go/src/github.com/pmokeev/rate-limiter

ENV GOPATH=/

COPY ./ ./
RUN apt-get update
RUN go mod download
RUN go build -o rate-limiter ./cmd/main.go

EXPOSE 8000
CMD ["./rate-limiter"]