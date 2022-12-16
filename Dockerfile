FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN make build

ENTRYPOINT ./main