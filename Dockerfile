FROM golang:1.13.4 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go get -u github.com/swaggo/swag/cmd/swag && swag init

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM alpine:3.10.3  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

COPY .env .

RUN mkdir -p ./data

COPY ./data/paradas.json ./data

COPY ./data/microbuses.json ./data

EXPOSE 8080

CMD ["./main"]