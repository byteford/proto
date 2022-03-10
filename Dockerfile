############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR /server/
COPY ./server .
# Fetch dependencies.
# Using go get.

RUN go get
RUN go build -o /go/bin/build
############################
# STEP 2 build a small image
############################
FROM alpine
# Copy our static executable.
COPY --from=builder /go/bin/build /go/bin/build
EXPOSE 3000
# Run the hello binary.
ENTRYPOINT ["/go/bin/build"]