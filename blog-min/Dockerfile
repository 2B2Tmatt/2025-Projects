FROM golang:1.24.5

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Bring in the rest of the source
COPY . .

# Build your app (if main.go is at repo root)
RUN go build -o server .

EXPOSE 9000
CMD ["./server"]
