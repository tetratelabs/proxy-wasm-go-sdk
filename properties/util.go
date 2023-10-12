package properties

import (
	"time"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// getIstioFilterMetadata parses istio filter metadata
func getIstioFilterMetadata(path []string) (IstioFilterMetadata, error) {
	result := IstioFilterMetadata{}

	config, err := getPropertyString(append(path, "config"))
	if err != nil {
		return IstioFilterMetadata{}, nil
	}
	result.Config = config

	services, err := getPropertyByteSliceSlice(append(path, "services"))
	if err != nil || services == nil {
		return result, nil
	}

	for _, service := range services {
		if service == nil {
			continue
		}
		istioService := IstioService{}
		istioServiceMap := deserializeStringMap(service)

		if host, ok := istioServiceMap["host"]; ok {
			istioService.Host = host
		}
		if name, ok := istioServiceMap["name"]; ok {
			istioService.Name = name
		}
		if namespace, ok := istioServiceMap["namespace"]; ok {
			istioService.Namespace = namespace
		}

		result.Services = append(result.Services, istioService)
	}

	return result, nil
}

// getPropertyBool returns a bool property.
func getPropertyBool(path []string) (bool, error) {
	bs, err := proxywasm.GetProperty(path)
	if err != nil {
		return false, err
	}

	return deserializeBool(bs)
}

// getPropertyByteSliceMap retrieves a complex property object as a map of byte slices.
// to be used when dealing with mixed type properties
func getPropertyByteSliceMap(path []string) (map[string][]byte, error) { //nolint:unused
	bs, err := proxywasm.GetProperty(path)
	if err != nil {
		return nil, err
	}

	return deserializeByteSliceMap(bs), nil
}

// getPropertyByteSliceSlice retrieves a complex property object as a string slice.
func getPropertyByteSliceSlice(path []string) ([][]byte, error) {
	bs, err := proxywasm.GetProperty(path)
	if err != nil {
		return nil, err
	}
	return deserializeByteSliceSlice(bs), nil
}

// getPropertyFloat64 returns a float64 property.
func getPropertyFloat64(path []string) (float64, error) {
	bs, err := proxywasm.GetProperty(path)
	if err != nil {
		return 0, err
	}

	return deserializeFloat64(bs), nil
}

// getPropertyString returns a string property.
func getPropertyString(path []string) (string, error) {
	bs, err := proxywasm.GetProperty(path)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

// getPropertyStringMap retrieves a complex property object as a map of string
// to be used when dealing with string only type properties.
func getPropertyStringMap(path []string) (map[string]string, error) {
	bs, err := proxywasm.GetProperty(path)
	if err != nil {
		return nil, err
	}

	return deserializeStringMap(bs), nil
}

// getPropertyStringSlice retrieves a  complex property object as a string slice.
func getPropertyStringSlice(path []string) ([]string, error) {
	bs, err := proxywasm.GetProperty(path)
	if err != nil {
		return nil, err
	}

	return deserializeStringSlice(bs), nil
}

// getPropertyTimestamp returns a timestamp property.
func getPropertyTimestamp(path []string) (time.Time, error) {
	bs, err := proxywasm.GetProperty(path)
	if err != nil {
		return time.Now().UTC(), err
	}

	return deserializeTimestamp(bs).UTC(), nil
}

// getPropertyUint64 returns a uint64 property.
func getPropertyUint64(path []string) (uint64, error) {
	bs, err := proxywasm.GetProperty(path)
	if err != nil {
		return 0, err
	}

	return deserializeUint64(bs), nil
}
