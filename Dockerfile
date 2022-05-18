FROM golang:1.17-buster

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go build -o main ./cmd/

CMD [ "./main" ]