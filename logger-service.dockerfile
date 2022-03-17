# The base go-image
FROM golang:1.18-alpine as builder

# Create a directory for the app
RUN mkdir /app

# Copy all files from the current directory to the app directory
COPY logger-service/. /app

# Set working directory
WORKDIR /app

# Run command as described:
# go build will build an executable file named server in the current directory
RUN CGO_ENABLED=0 go build -o logServiceApp ./cmd/web

RUN chmod +x /app/logServiceApp

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app
RUN mkdir /templates

COPY --from=builder /app/logServiceApp /app
COPY --from=builder /app/templates/. /templates

# Run the server executable
CMD [ "/app/logServiceApp" ]