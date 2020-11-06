FROM golang:1.15.3-alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN swagger 

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o app ./cmd/web

# Export necessary port
EXPOSE 4000

# Command to run when starting the container
CMD ["/build/app"]