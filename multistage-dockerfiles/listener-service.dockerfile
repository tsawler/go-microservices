# The base go-image
FROM golang:1.18-alpine as builder
 
# create a directory for the app
RUN mkdir /app
 
# copy all files from the current directory to the app directory
COPY listener-service /app
 
# set working directory
WORKDIR /app
 
# build executable
RUN CGO_ENABLED=0 go build -o listener .

RUN chmod +x /app/listener

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app

COPY --from=builder /app/listener /app

# Run the server executable
CMD [ "/app/listener" ]