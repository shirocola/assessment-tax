# Start by building the application.
FROM golang:1.21 as builder

WORKDIR /go/src/app
COPY . .

# Use go mod download to pull in any dependencies
RUN go mod download

# Build your application on the builder image.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10

COPY --from=builder /go/src/app/server /server

# Run the server binary.
ENTRYPOINT ["/server"]
