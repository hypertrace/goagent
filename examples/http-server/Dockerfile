# docker build -t tmwangi/hackgoapp:0.1.0 .
FROM ubuntu:20.10

WORKDIR /app
COPY http-server http-server
COPY test_config.yaml test_config.yaml

ENV HT_CONFIG_FILE test_config.yaml

ENTRYPOINT ["./http-server"]
