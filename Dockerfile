FROM golang:1.25.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY Project1.go ./

RUN go build -o xkcd-server Project1.go

EXPOSE 8080

CMD [ "./xkcd-server", "--server" ]