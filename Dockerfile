FROM golang:1.12.8 AS builder
WORKDIR /app
COPY rsimkin /app/rsimkin/
ADD *.go /app/
ADD go.* /app/
RUN go build -o /app/app
COPY config /app/config/

FROM golang:1.12.8
WORKDIR /app
COPY --from=builder /app/app /app
COPY --from=builder /app/config/ /app/config/
CMD ["/app/app"]