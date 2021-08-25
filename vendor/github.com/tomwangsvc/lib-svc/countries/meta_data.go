package countries

import (
	"context"
	"encoding/json"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
)

// NewMetadata is a factory for Metadata interfaces
func NewMetadata(ctx context.Context) (*Metadata, error) {
	lib_log.Info(ctx, "Initializing")

	countries, err := loadMetadata()
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed to load meta-data")
	}

	var countriesByCountryCode = make(map[string]Country)
	var countriesByCurrencyCode = make(map[string][]Country)

	for _, v := range countries {
		if _, ok := countriesByCountryCode[v.CountryCode]; !ok {
			countriesByCountryCode[v.CountryCode] = v
		}

		countriesByCurrencyCode[v.CurrencyCode] = append(countriesByCurrencyCode[v.CurrencyCode], v)
	}

	return &Metadata{countriesByCountryCode, countriesByCurrencyCode}, nil
}

// Metadata is an interface for working with country meta-data
type Metadata struct {
	CountriesByCountryCode  map[string]Country
	CountriesByCurrencyCode map[string][]Country
}

func loadMetadata() ([]Country, error) {
	var c []Country
	err := json.Unmarshal([]byte(jsonMetadata()), &c)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed to parse embedded countries")
	}
	return c, nil
}
