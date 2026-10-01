FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY carsAPI.go .

RUN go build -o carsAPI carsAPI.go

FROM scratch

COPY --from=builder /app/carsAPI /carsAPI

EXPOSE 3015
ENTRYPOINT ["/carsAPI"]
