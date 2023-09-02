FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o ./portfolio

EXPOSE 8000

CMD ["./portfolio"]