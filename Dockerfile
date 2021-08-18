FROM golang:1.16-alpine as builder
WORKDIR /klein

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o klein

FROM golang:1.16-alpine

WORKDIR /klein

COPY --from=builder /klein/klein ./
COPY --from=builder /klein/config/config.yaml ./config/
COPY --from=builder /klein/views /klein/views

EXPOSE 8080

CMD ["./klein"]

