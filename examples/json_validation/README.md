# JSON Payload validation

This wasm plugin checks whether the request has JSON payload and has required keys in it.
If not, the wasm plugin ceases the further process of the request and returns 403 immediately.

## Run it via Envoy

`envoy.yaml` is the example envoy config file that you can use for running the wasm plugin
with standalone Envoy.

Envoy listens on `localhost:18000`, responding to any requests with static content "hello from server".
However, the wasm plugin also runs to validate the requests' payload.

```bash
make -C ../.. run name=json_validation
```

The plugin intercepts the request and makes Envoy return 403 instead of the static content
if the request has no JSON payload or the payload JSON doesn't have "id" or "token" keys.

```console
# Returns the normal response when the request has the required keys, id and token.
curl -X POST localhost:18000 -H 'Content-Type: application/json' --data '{"id": "xxx", "token": "xxx"}'
hello from the server

# Returns 403 when the request has missing required keys.
curl -v -X POST localhost:18000 -H 'Content-Type: application/json' --data '"required_keys_missing"'
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:18000...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 18000 (#0)
> POST / HTTP/1.1
> Host: localhost:18000
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 23
>
* upload completely sent off: 23 out of 23 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 403 Forbidden
< content-length: 15
< content-type: text/plain
< date: Tue, 01 Mar 2022 19:22:24 GMT
< server: envoy
<
* Connection #0 to host localhost left intact
invalid payload
```

### Run it via Istio

This example details deploying to a kind cluster running the Istio httpbin sample app.

```console
# Create a test cluster
kind create cluster

# Install Istio and the httpbin sample app

istioctl install --set profile=demo -y
kubectl label namespace default istio-injection=enabled
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.12/samples/httpbin/httpbin.yaml
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.12/samples/httpbin/httpbin-gateway.yaml
```

For Istio 1.12 and later the easiest way is to use a WasmPlugin resource. For older Istio
versions an EnvoyFilter is needed.

#### Install using WasmPlugin resource

Build and push the wasm module to your container registry, then apply the WasmPlugin.

```console
export HUB=your_registry # e.g. docker.io/tetrate
make -C ../.. build.example name=json_validation
docker build . -t ${HUB}/json-validation:v1
docker push ${HUB}/json-validation:v1

sed "s|YOUR_CONTAINER_REGISTRY|$HUB|" wasmplugin.yaml | kubectl apply -f -
```

#### Install using EnvoyFilter

To use an EnvoyFilter, create a config map containing the compiled wasm plugin, mount the config
map into the gateway pod, and then configure Envoy via an EnvoyFilter to load the wasm plugin from
a local file.

```console
# Create the config map
kubectl -n istio-system create configmap wasm-plugins --from-file=main.wasm

# Patch the gateway deployment to mount the config map
kubectl -n istio-system patch deployment istio-ingressgateway --patch-file=gatewaydeploymentpatch.yaml

# Create the EnvoyFilter
kubectl apply -f envoyfilter.yaml
```

#### Test the plugin

Expose the ingress gateway on port 8080 on your local machine via.

```console
kubectl port-forward -n istio-system svc/istio-ingressgateway 8080:80
```

Requests without the required payload will fail:

```console
% curl -i http://localhost:8080/post  -H 'Content-Type: application/json' --data '{"id": "xxx", "not_token": "xxx"}'
HTTP/1.1 403 Forbidden
content-length: 15
content-type: text/plain
date: Tue, 15 Mar 2022 23:05:52 GMT
server: istio-envoy

invalid payload
```

But those with the payload will proceed:

```console
% curl -i http://localhost:8080/post  -H 'Content-Type: application/json' --data '{"id": "xxx", "token": "xxx"}'
HTTP/1.1 200 OK
server: istio-envoy
date: Tue, 15 Mar 2022 23:06:29 GMT
content-type: application/json
content-length: 884
access-control-allow-origin: *
access-control-allow-credentials: true
x-envoy-upstream-service-time: 3

{
  "args": {},
  "data": "{\"id\": \"xxx\", \"token\": \"xxx\"}",
  "files": {},
  "form": {},
  "headers": {
    "Accept": "*/*",
    "Content-Length": "29",
    "Content-Type": "application/json",
    "Host": "localhost:8080",
    "User-Agent": "curl/7.64.1",
    "X-B3-Parentspanid": "99a94908edd26592",
    "X-B3-Sampled": "1",
    "X-B3-Spanid": "e12fc7fd9aa74838",
    "X-B3-Traceid": "2b7375cda8bc98a299a94908edd26592",
    "X-Envoy-Attempt-Count": "1",
    "X-Envoy-Internal": "true",
    "X-Forwarded-Client-Cert": "By=spiffe://cluster.local/ns/default/sa/httpbin;Hash=5703a66dcdbc8cafc8c29e1ebee1174f4bc81234d8dc1ccc20fb9e3c26b320e1;Subject=\"\";URI=spiffe://cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account"
  },
  "json": {
    "id": "xxx",
    "token": "xxx"
  },
  "origin": "10.244.0.9",
  "url": "http://localhost:8080/post"
}
```
