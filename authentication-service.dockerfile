# The base go-image
FROM golang:1.18-alpine as builder

# Create a directory for the app
RUN mkdir /app

# Copy all files from the current directory to the app directory
COPY authentication-service/. /app

# Set working directory
WORKDIR /app

# Run command as described:
# go build will build an executable file named server in the current directory
RUN CGO_ENABLED=0 go build -o authApp .

RUN chmod +x /app/authApp

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app

COPY --from=builder /app/authApp /app

# Run the server executable
CMD [ "/app/authApp" ]