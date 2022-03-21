FROM alpine:latest
RUN mkdir /app
RUN mkdir /templates

COPY logger-service/logServiceApp /app
COPY logger-service/templates/. /templates

# Run the server executable
CMD [ "/app/logServiceApp" ]