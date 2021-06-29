FROM envoyproxy/envoy:v1.18.3 AS envoy
WORKDIR /app
RUN cp $(which envoy) ./envoy

FROM golang:1.16 AS go

FROM tinygo/tinygo:0.18.0
ENV PATH $PATH:/usr/local/go/bin
RUN apt-get clean && \
    apt-get update && \
    apt-get install -y make
COPY --from=envoy /app/envoy /usr/local/bin
COPY --from=go /usr/local/go /usr/local/go
