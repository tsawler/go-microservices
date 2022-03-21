# The base go-image
FROM golang:1.18-alpine as builder

# create a directory for the app
RUN mkdir /app

# copy all files from the current directory to the app directory
COPY mail-service/. /app

# set working directory
WORKDIR /app

# build executable
RUN CGO_ENABLED=0 go build -o mailerServiceApp .

RUN chmod +x /app/mailerServiceApp

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app
RUN mkdir /templates

COPY --from=builder /app/mailerServiceApp /app
COPY --from=builder /app/templates /templates

# Run the server executable
CMD [ "/app/mailerServiceApp" ]