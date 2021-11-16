# Istio Wasm Demo for Rate Limiting

Istio 1.12 release introduces new Wasm Extension API. This folder contains a sample application of
implementing rate limiting in Golang, and deploy the Wasm Plugin using Istio API.

<!-- TODO(incfly): provide the actual link once the blog is ready. -->
For detailed instructions, checkout tetrate.io/blog.

##  Build and publish Wasm extension on OCI registry


Dependency:: TinyGo, and Docker CLI


1. Compile the code to Wasm binary.

```
$ tinygo build -o main.wasm -scheduler=none -target=wasi main.go
```

2. Build docker image which is compliant with Wasm OCI image spec (https://github.com/solo-io/wasm/tree/master/spec)

```
# Note that replace ${WASM_EXTENSION_REGISTRY} with your OCI repo.
# Here I push to GitHub Container Registry.
$ export WASM_EXTENSION_REGISTRY=ghcr.io/mathetake/wasm-extension-demo
$ docker build -t ${WASM_EXTENSION_REGISTRY}:v1 .
```

3. Publish the docker image to your OCI registry

```
# Make sure you already logged in to the registory.
docker push ${WASM_EXTENSION_REGISTRY}:v1
```

