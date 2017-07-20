FROM debian:jessie
RUN apt-get update
RUN apt-get install -y ca-certificates
COPY ./linked-air /linked-air
ENTRYPOINT ["/linked-air"]
