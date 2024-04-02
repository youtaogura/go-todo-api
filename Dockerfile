FROM golang:1.22.1

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.29.0
RUN apt-get update && apt-get install -y netcat-openbsd

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY boot.sh ./
CMD ["sh", "boot.sh"]
