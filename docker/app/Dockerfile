# snippetbox
FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN mkdir -p tls \
    && cd tls && go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost \
    && cd ..

COPY . .

RUN go get ./pkg/forms \
    && go get ./pkg/models \
    && go get ./pkg/models/mysql

EXPOSE 4000

RUN go build -o ./app ./cmd/web

CMD [ "./app" ]