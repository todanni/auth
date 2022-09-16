# build stage
FROM golang:alpine as build

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY . /build

WORKDIR /build

# Build the application
RUN go build -o main .


# actual image
FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
COPY --from=build /build/main .

# Export necessary port
EXPOSE 8083

# Command to run when starting the container
CMD ["/dist/main"]
