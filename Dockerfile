# FROM golang:1.17

# WORKDIR $HOME/go/src/X-Blog

# COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

# COPY . ./X-Blog

 
# EXPOSE 8081



FROM golang:1.17-buster

#ENV TZ=Europe/Moscow
#RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone


WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go build -o main ./cmd/

CMD [ "./main" ]

# ENTRYPOINT CompileDaemon --build="go build cmd/main.go" --command=./main