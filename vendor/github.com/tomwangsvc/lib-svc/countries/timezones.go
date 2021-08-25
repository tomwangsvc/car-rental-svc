package countries

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_http "github.com/tomwangsvc/lib-svc/http"
	lib_log "github.com/tomwangsvc/lib-svc/log"
)

type GeoNamesTimezone struct {
	CountryName string  `json:"countryName"` // ISO 3166 country code name.
	CountryCode string  `json:"countryCode"` // ISO 3166 country code.
	DstOffset   float64 `json:"dstOffset"`   // Offset to GMT at 1. July (deprecated).
	GmtOffset   float64 `json:"gmtOffset"`   // Offset to GMT at 1. January (deprecated).
	Lat         float64 `json:"lat"`         // Lat used for the call.
	Lng         float64 `json:"lng"`         // Lng used for the call.
	RawOffset   float64 `json:"rawOffset"`   // The amount of time in hours to add to UTC to get standard time in this time zone.
	Sunrise     string  `json:"sunrise"`     // Current days time for sunrise.
	Sunset      string  `json:"sunset"`      // Current days time for sunset.
	Time        string  `json:"time"`        // The local current time.
	TimezoneID  string  `json:"timezoneID"`  // The name of the timezone (according to olson).
}

func ReadTimezone(ctx context.Context, lat, lng float64) (*GeoNamesTimezone, error) {
	lib_log.Info(ctx, "Reading", lib_log.FmtFloat64("lat", lat), lib_log.FmtFloat64("lng", lng))

	client := lib_http.NewClient(ctx, true)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://api.geonames.org/timezoneJSON?lat=%f&lng=%f&username=%s", lat, lng, GeoNamesUsername), nil)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed creating searching geonames timezones request")
	}
	req = req.WithContext(ctx)
	req.Header.Add("Accept", "application/json")

	lib_http.LogRequest(req)

	res, err := client.Do(req)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed issuing geonames timezones request")
	}
	defer lib_http.CloseBody(ctx, res.Body)

	body, err := lib_http.ReadResponseBody(res, client, false)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed reading geonames timezones response body")
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, lib_errors.NewCustomf(http.StatusNotFound, "Geonames timezones not found, request resulted in http status code %d", res.StatusCode)
	} else if res.StatusCode != http.StatusOK {
		return nil, lib_errors.NewCustomf(http.StatusBadGateway, "Failed geonames timezones request, request resulted in http status code %d", res.StatusCode)
	}

	var geoNamesTimezone GeoNamesTimezone
	if err := json.Unmarshal(body, &geoNamesTimezone); err != nil {
		return nil, lib_errors.Wrap(err, "Failed unmarshalling response body into geoNamesTimezone")
	}

	if geoNamesTimezone.TimezoneID == "" {
		return nil, lib_errors.NewCustomf(http.StatusNotFound, "Invalid geonames timezone found, missing timezoneID: %v", geoNamesTimezone)
	}

	lib_log.Info(ctx, "Read", lib_log.FmtAny("geoNamesTimezone", geoNamesTimezone))
	return &geoNamesTimezone, nil
}
