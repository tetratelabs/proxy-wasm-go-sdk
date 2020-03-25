package runtime

const intSize = 32 << (^uint(0) >> 63)

func init() {
	// just in case, not sure this is necessary:)
	if intSize != 32 {
		panic("must be 32bit")
	}
}

//go:export proxy_abi_version_0_1_0
func proxyABIVersion010() {}
