# The base go-image
FROM golang:1.18-alpine as builder

# create a directory for the app
RUN mkdir /app

# copy all files from the current directory to the app directory
COPY authentication-service/. /app

# set working directory
WORKDIR /app

# build executable
RUN CGO_ENABLED=0 go build -o authApp .

RUN chmod +x /app/authApp

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app

COPY --from=builder /app/authApp /app

# Run the server executable
CMD [ "/app/authApp" ]