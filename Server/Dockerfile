FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod init chatbox || true
RUN go mod tidy

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]