FROM golang:alpine

WORKDIR /usr/src/app
COPY *.go go.mod go.sum game.html key.bin ./
RUN go build -v

FROM httpd:alpine

WORKDIR /app
COPY --from=0 /usr/src/app/cgi /usr/local/apache2/htdocs/index.cgi

WORKDIR /usr/local/apache2/htdocs
COPY index.css ./
