FROM alpine:latest
RUN mkdir /app
RUN mkdir /templates

COPY mail-service/mailerServiceApp /app
COPY mail-service/templates /templates

# Run the server executable
CMD [ "/app/mailerServiceApp" ]