FROM golang:1.15.3-buster AS builder

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

# Copy the code into the container
ADD ./swag-1.6.9.tar.gz ./
COPY . .
# Build the application

RUN ls
RUN ./swag init -g ./cmd/web/main.go
RUN go build -o app ./cmd/web

# Export necessary port
EXPOSE 4000

FROM scratch

COPY --from=builder /build/app /
COPY --from=builder /build/docs/* /

CMD ["/app"]