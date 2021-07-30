# Dockerfile for building "compat" variant of Wasm Image Specification.
# https://github.com/solo-io/wasm/blob/master/spec/spec-compat.md
FROM scratch

ARG WASM_BINARY_PATH
COPY ${WASM_BINARY_PATH} ./plugin.wasm
