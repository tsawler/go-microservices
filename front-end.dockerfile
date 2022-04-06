FROM alpine:latest
RUN mkdir /app

COPY front-end/frontEndLinux /app

# Run the server executable
CMD [ "/app/frontEndLinux" ]