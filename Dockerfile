# consignment-service/Dockerfile

# Use officel golang image
# Refer to the image as `builder`.
FROM golang:1.9.0 as builder

# Set workdir to current service in gopath.
WORKDIR /go/src/github.com/Terry-Bui/go-microservice-example/consignment-service

# Copy current code into workdir.
COPY . .


# Pull in the dependency management tool godep.
RUN go get -u github.com/golang/dep/cmd/dep

# Create a dep project, run `ensure`, which will pull in
# all of the depenndencies within this directory.
RUN dep init & dep ensure

# Build binary & run the binary in Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

# Tell Docker to start new build process with alpine image
FROM alpine:latest

# Add central authority certificates
RUN apk --no-cache add ca-certificates

# Create directory for the app
RUN mkdir /app
WORKDIR /app

# Pull binary from the `builder` container within the build context.
# This gets the previous image, find the binary that was built and
# pulls it into this container
COPY --from=builder /go/src/Terry-Bui/consignment-service/consignment-service .

# Run binary and build it in a seperate container with all of the
# dependencies and run time libraries.
CMD ["./consignment-service"]

