FROM golang:1.20.2-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./

RUN go build -o vault-pract .

FROM alpine:3.17.2

COPY --from=builder /app/vault-pract .
EXPOSE 8181
CMD [ "./vault-pract" ]