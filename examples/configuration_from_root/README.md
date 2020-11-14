## configuration_from_root

This example reads the json string from Envoy's configuration yaml at the startup time
The child HTTP context then reads the config from its corresponding root context.

```
wasm log my_root_id: plugin config: {
  "name": "plugin configuration"
}
```
