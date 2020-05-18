package commercio

// We use init() to configure commercionetwork's Cosmos settings, otherwise address/codec won't work
func init() {
	setCosmosConfig()
}
