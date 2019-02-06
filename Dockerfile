FROM debian:stable-slim

WORKDIR /opt/poc/bin
COPY bin .

EXPOSE 5555

CMD ["/opt/poc/bin/ipblock-pool-server"]