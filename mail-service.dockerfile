# The base go-image
FROM golang:1.17-alpine

# Create a directory for the app
RUN mkdir /app

# Copy all files from the current directory to the app directory
COPY mail-service/. /app

# Set working directory
WORKDIR /app

# Run command as described:
# go build will build an executable file named server in the current directory
RUN CGO_ENABLED=0 go build -o mailerServiceApp .

RUN chmod +x /app/mailerServiceApp

# Run the server executable
CMD [ "/app/mailerServiceApp" ]