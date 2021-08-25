package countries

var (
	// BlockListCountries are countries not served by tomwang
	BlockListCountries = make(map[string]interface{})
)

func init() {
	BlockListCountries["KP"] = "KP"
	BlockListCountries["CU"] = "CU"
}

// IsBlockListCountry check if the country is in the blocklist
func IsBlockListCountry(countryCode string) bool {
	if _, ok := BlockListCountries[countryCode]; ok {
		return true
	}
	return false
}
