# The base go-image
FROM alpine:latest
RUN mkdir /app

COPY broker-service/brokerApp /app

# Run the server executable
CMD [ "/app/brokerApp" ]