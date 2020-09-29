FROM istio/proxyv2:1.7.2

FROM golang:1.15

COPY --from=0 /usr/local/bin/envoy /usr/local/bin/
