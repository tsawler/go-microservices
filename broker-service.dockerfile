# The base go-image
FROM golang:1.18-alpine as builder

# create a directory for the app
RUN mkdir /app

# copy all files from the current directory to the app directory
COPY broker-service/. /app

# set working directory
WORKDIR /app

# build executable
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app

COPY --from=builder /app/brokerApp /app

# Run the server executable
CMD [ "/app/brokerApp" ]