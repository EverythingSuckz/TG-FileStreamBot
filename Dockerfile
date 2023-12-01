FROM golang:1.21 AS builder
RUN apt-get update && apt-get upgrade -y && apt-get install build-essential -y
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build ./cmd/fsb/ -o out/fsb -ldflags="-w -s" .

FROM golang:1.21
COPY --from=builder /app/out/fsb /app/fsb
CMD ["/app/fsb"]