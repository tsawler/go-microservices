FROM alpine:latest
RUN mkdir /app

COPY authentication-service/authApp /app

# Run the server executable
CMD [ "/app/authApp" ]