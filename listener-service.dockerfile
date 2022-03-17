# The base go-image
FROM golang:1.18-alpine as builder
 
# Create a directory for the app
RUN mkdir /app
 
# Copy all files from the current directory to the app directory
COPY listener-service /app
 
# Set working directory
WORKDIR /app
 
# Run command as described:
# go build will build an executable file named server in the current directory
RUN CGO_ENABLED=0 go build -o listener .

RUN chmod +x /app/listener

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app

COPY --from=builder /app/listener /app

# Run the server executable
CMD [ "/app/listener", "log.INFO", "log.WARNING", "log.ERROR" ]