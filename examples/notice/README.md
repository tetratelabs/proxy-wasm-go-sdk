## notice
`notice` 通过使用`EnvoyFilter`和`wasm`，达到在`bookinfo`的`productPage`服务中发布公告的功能。

``EnvoyFilter``的配置基于envoy的静态配置形成，配置起来很复杂，这个以后再花时间去好好深入。配置``EnvoyFilter``有以下几个注意点：

- `proxyVersion` 需要与 `istio-proxy` 版本保持一致。我使用的是`istio` 1.13，所以 `proxyVersion: ^1\.13.*`。
-  当wasm使用local的时候，`filename` 需要与前面`Deployment`中的 Annotation 保持一致。
- `workloadSelector` 设置为目标 `Pod` 的 label。

