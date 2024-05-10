FROM golang:latest

WORKDIR /app

COPY . .

RUN go build

CMD ["./yadrotest", "test_file.txt"]