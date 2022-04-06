FROM alpine:latest
RUN mkdir /app

COPY front-end/frontEndLinux /app
COPY front-end/templates /app/templates

WORKDIR /app

# Run the server executable
CMD [ "./frontEndLinux" ]