FROM golang:1.23-alpine3.20

WORKDIR /app

COPY go.mod .

RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq

COPY . .

RUN go build -o github.com/scottbrazil/golang-postgres-api-study-docker .

EXPOSE 7777

CMD ["./github.com/scottbrazil/golang-postgres-api-study-docker"]