FROM golang:1.18-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /love

FROM alpine
COPY --from=builder /love /love
EXPOSE 8080
ENTRYPOINT [ "/love" ]
