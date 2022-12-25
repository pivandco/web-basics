FROM golang:1.19.4-alpine

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .
RUN go build -o backend .

CMD ["/app/backend"]
