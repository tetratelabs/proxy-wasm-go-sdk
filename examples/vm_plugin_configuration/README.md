## vm_plugin_configuration

This example reads the json string from Envoy's configuration yaml at the startup time


```
wasm log my_root_id: vm config: {
  "name": "vm configuration"
}


wasm log my_root_id: plugin config: {
  "name": "plugin configuration"
}


wasm log my_root_id: vm config: {
  "name": "vm configuration"
}


wasm log my_root_id: plugin config: {
  "name": "plugin configuration"
}
```
