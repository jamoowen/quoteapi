# STAGE 1: Builder
FROM golang:1.24-alpine AS builder    

WORKDIR /app                         

RUN apk add --no-cache gcc musl-dev       

# Copy go mod and sum files
COPY go.mod go.sum ./               

# Download dependencies
RUN go mod download                 

# Copy the source code
COPY . .                           

# Build the application

RUN CGO_ENABLED=1 GOOS=linux go build -o bin/quoteapi ./cmd/server/main.go

# STAGE 2: Final
FROM alpine:latest                 

# Install required dependencies for SQLite
RUN apk add --no-cache sqlite libc6-compat    

WORKDIR /app                      

RUN mkdir -p /app/bin /app/db /app/data/csv 

# Copy only the binary from builder
COPY --from=builder /app/bin/quoteapi /app/bin/quoteapi   

# Debug: Print the contents of the source directory
RUN echo "copying static data..." 
# csv
COPY data/quotes.csv /app/data/
# the sql migration
COPY migrations /app/migrations
COPY templates /app/templates

# Run the binary
CMD ["/app/bin/quoteapi"]                    
