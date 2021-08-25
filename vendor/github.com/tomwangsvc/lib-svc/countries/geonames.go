package countries

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_http "github.com/tomwangsvc/lib-svc/http"
	lib_log "github.com/tomwangsvc/lib-svc/log"
)

const (
	defaultLanguage = "en"
)

var (
	countryCache     map[string][]Country
	geosCache        map[string][]Geo
	countryCacheTime time.Time
)

func init() {
	countryCache = make(map[string][]Country)
	geosCache = make(map[string][]Geo)
}

// ReadCountryMetadataFromGeonames reads country meta-data from GeoNames
func ReadCountryMetadataFromGeonames(ctx context.Context, language, countryCode string) ([]Country, error) {
	lib_log.Debug(ctx, "Enter")
	defer lib_log.Debug(ctx, "Leave")

	if v, ok := countryCache[language]; ok {
		for _, vv := range v {
			if vv.CountryCode == countryCode {
				return []Country{vv}, nil
			}
		}
	}

	c, err := readCountriesFromGeonames(ctx, language, countryCode)
	if err != nil {
		return nil, err
	}
	if v, ok := countryCache[language]; ok {
		var found bool
		for _, vv := range v {
			if vv.CountryCode == countryCode {
				found = true
				break
			}
		}
		if !found {
			for _, country := range c {
				if country.CountryCode == countryCode {
					countryCache[language] = append(countryCache[language], country)
				}
			}
		}
	}

	return c, nil
}

// ReadAllCountryMetadataFromGeonames reads country meta-data from GeoNames for all countries
func ReadAllCountryMetadataFromGeonames(ctx context.Context, language string) ([]Country, error) {
	lib_log.Debug(ctx, "Enter")
	defer lib_log.Debug(ctx, "Leave")

	now := time.Now()
	if !countryCacheTime.Before(now.Add(-time.Hour * 24)) {
		if v, ok := countryCache[language]; ok {
			if len(v) > 0 {
				return v, nil
			}
		}
	}

	c, err := readCountriesFromGeonames(ctx, language, "")
	if err != nil {
		lib_log.Warn(ctx, "Failed reading countries from geonames, will try to fall back to cache", lib_log.FmtError(err))
		if v, ok := countryCache[language]; ok {
			if len(v) > 0 {
				return v, nil
			}
		} else if language != defaultLanguage {
			if v, ok := countryCache[defaultLanguage]; ok {
				if len(v) > 0 {
					return v, nil
				}
			}
		}
		return nil, lib_errors.Wrap(err, "Failed reading countries from geonames")
	}

	countryCache[language] = c
	countryCacheTime = now

	return countryCache[language], nil
}

// ReadAllGeoMetadataFromGeonames reads geo meta-data from GeoNames
func ReadAllGeoMetadataFromGeonames(ctx context.Context, language, geonameID, criteria string) ([]Geo, error) {
	lib_log.Debug(ctx, "Enter")
	defer lib_log.Debug(ctx, "Leave")

	if criteria != GeoNameCriteriaNeighbours && criteria != GeoNameCriteriaChildren && criteria != GeoNameCriteriaSiblings && criteria != GeoNameCriteriaContains && criteria != GeoNameCriteriaHierarchy {
		return nil, lib_errors.NewCustomf(http.StatusBadRequest, "Not recognized: criteria=%s, expected %q or %q or %q or %q or %q", criteria, GeoNameCriteriaNeighbours, GeoNameCriteriaChildren, GeoNameCriteriaSiblings, GeoNameCriteriaContains, GeoNameCriteriaHierarchy)
	}

	lib_log.Info(ctx, "Read geonameID", lib_log.FmtString("geonameID", geonameID))

	if v, ok := geosCache[language+geonameID+criteria]; ok {
		if len(v) > 0 {
			return v, nil
		}
	}

	c, err := searchGeos(ctx, language, geonameID, criteria)
	if err != nil {
		return nil, err
	}

	geosCache[language+geonameID+criteria] = c

	return geosCache[language+geonameID+criteria], nil
}

func readCountriesFromGeonames(ctx context.Context, language, countryCode string) ([]Country, error) {
	lib_log.Debug(ctx, "Enter")
	defer lib_log.Debug(ctx, "Leave")

	client := lib_http.NewClient(ctx, true)
	req, err := http.NewRequest(http.MethodGet, "http://api.geonames.org/countryInfoJSON?style=full", nil)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed creating GeoNames read countries request")
	}
	req = req.WithContext(ctx)

	q := req.URL.Query()
	q.Set("lang", language)
	q.Set("username", GeoNamesUsername)
	if countryCode != "" {
		q.Set("country", countryCode)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")

	lib_http.LogRequest(req)

	res, err := client.Do(req)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed issuing GeoNames read countries request")
	}
	defer lib_http.CloseBody(ctx, res.Body)

	body, err := lib_http.ReadResponseBody(res, client, false)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed reading GeoNames read countries response body")
	}

	if res.StatusCode != http.StatusOK {
		return nil, lib_errors.NewCustomf(http.StatusBadGateway, "Failed GeoNames read countries request: Request resulted in http status code %d", res.StatusCode)
	}

	var geonamesCountryRes GeonamesCountryRes
	if err := json.Unmarshal(body, &geonamesCountryRes); err != nil {
		return nil, lib_errors.Wrap(err, "Failed unmarshalling GeoNames read countries response body into geonamesCountryRes")
	}

	lib_log.Debug(ctx, "Received GeoNames countries", lib_log.FmtString("countryCode", countryCode), lib_log.FmtString("language", language), lib_log.FmtString("countryCode", countryCode), lib_log.FmtAny("len(geonamesCountryRes.Geonames)", len(geonamesCountryRes.Geonames)))

	return NewCountriesFromGeonamesCountry(ctx, geonamesCountryRes.Geonames), nil
}

func searchGeos(ctx context.Context, language, geonameID, criteria string) ([]Geo, error) {
	lib_log.Debug(ctx, "Enter")
	defer lib_log.Debug(ctx, "Leave")

	if geonameID == "" {
		return nil, lib_errors.NewCustom(http.StatusBadRequest, "Failed attempting read GeoNames geos request: Missing geonameID")
	}

	if criteria == "" {
		return nil, lib_errors.NewCustom(http.StatusBadRequest, "Failed attempting read GeoNames geos request: Missing criteria")
	}

	url := fmt.Sprintf("http://api.geonames.org/%sJSON?style=full", criteria)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed reading GeoNames geos request")
	}
	req = req.WithContext(ctx)

	q := req.URL.Query()
	q.Set("lang", language)
	q.Set("geonameId", geonameID)
	q.Set("username", GeoNamesUsername)

	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")

	lib_http.LogRequest(req)

	client := lib_http.NewClient(ctx, true)
	res, err := client.Do(req)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed issuing GeoNames read geos request")
	}
	defer lib_http.CloseBody(ctx, res.Body)

	body, err := lib_http.ReadResponseBody(res, client, false)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed reading GeoNames read geos response body")
	}

	if res.StatusCode != http.StatusOK {
		return nil, lib_errors.NewCustomf(http.StatusBadGateway, "Failed reading GeoNames geos request: Request resulted in http status code %d", res.StatusCode)
	}

	var geonamesGeoRes GeonamesGeoRes
	if err := json.Unmarshal(body, &geonamesGeoRes); err != nil {
		return nil, lib_errors.Wrap(err, "Failed unmarshalling GeoNames read geos response body into geonamesGeoRes")
	}

	lib_log.Debug(ctx, "Received GeoNames geos", lib_log.FmtString("geonameID", geonameID), lib_log.FmtInt("geonamesGeoRes.TotalResultsCount", geonamesGeoRes.TotalResultsCount), lib_log.FmtInt("len(geonamesGeoRes.Geonames)", len(geonamesGeoRes.Geonames)))

	geos, err := NewGeosFromGeonamesGeos(geonamesGeoRes.Geonames)
	if err != nil {
		return nil, err
	}
	return geos, nil
}
