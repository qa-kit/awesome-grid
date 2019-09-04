FROM golang:1.12.8 AS builder
WORKDIR /app
COPY ./ /app/
RUN go build -o /app/app

FROM golang:1.12.8
WORKDIR /app
COPY --from=builder /app/app /app
COPY --from=builder /app/*.json /app/
CMD ["/app/app"]