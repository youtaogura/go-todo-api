FROM golang:1.22.1

RUN go install github.com/cosmtrek/air@v1.29.0

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
