FROM golang:1.24-alpine AS builder

WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=1
RUN go build -o carsAPI .

FROM alpine
WORKDIR /app
COPY --from=builder /app/carsAPI /app/carsAPI

EXPOSE 3015
ENTRYPOINT ["/app/carsAPI"]
