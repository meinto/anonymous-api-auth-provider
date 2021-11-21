FROM golang:1.17-alpine AS builder
WORKDIR /service
COPY . .
RUN go build -o authprovider .

FROM alpine
RUN apk add bash
WORKDIR /service
COPY --from=builder /service/authprovider .
CMD [ "/service/authprovider" ]