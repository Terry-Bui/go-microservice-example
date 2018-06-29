FROM debian:latest

RUN mkdir -p /app
WORKDIR /app

ADD consignment.json /app/consignment.json
ADD consignment-client /app/consignment-client

CMD ["./consignment-client"]

