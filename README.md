> [!WARNING] 
> We are no longer recommending this SDK or Wasm in general for anyone due to the fundamental memory issue of TinyGo (See [the detailed explanation](https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/450#issuecomment-2253729297) by a long-time community member)
> as well as [the project state of Proxy-Wasm in general](https://github.com/envoyproxy/envoy/issues/35420).
> If you are not in a position where you have to run untrusted binaries (like for example, you run Envoy proxies while your client gives you the binaries to run), we recommend using other extension mechanism 
> such as Lua or External Processing which should be comparable or better or worse depending on the use case. 
> 
> If you are already using this SDK, but still want to continue using Wasm for some reason instead of Lua or External Processing, 
> we strongly recommend migrating to the Rust or C++ SDK due to the memory issue of TinyGo described in the link above.
> 
> We keep this repository open and not archived for the existing users, but we cannot provide any support or guarantee for the future development of this SDK.
> However, at any time, we may decide to archive this repository if we see no reason to keep it open.
