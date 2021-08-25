package countries

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	"github.com/ttacon/libphonenumber"
)

// GeonamesCountryRes data model
type GeonamesCountryRes struct {
	Geonames []GeonamesCountry `json:"geonames"`
}

// GeonamesCountry data model
/*
    {
      "continent": "AS",
      "capital": "Singapore",
      "languages": "cmn,en-SG,ms-SG,ta-SG,zh-SG",
      "geonameId": 1880251,
      "south": 1.258556,
      "isoAlpha3": "SGP",
      "north": 1.471278,
      "fipsCode": "SN",
      "population": "4701069",
      "east": 104.007469,
      "isoNumeric": "702",
      "areaInSqKm": "692.7",
      "countryCode": "SG",
      "west": 103.638275,
      "countryName": "Singapore",
      "continentName": "Asia",
      "currencyCode": "SGD"
		}
*/
type GeonamesCountry struct {
	AreaInSqKm    string `json:"areaInSqKm"`
	Capital       string `json:"capital"`
	Continent     string `json:"continent"`
	ContinentName string `json:"continentName"`
	CountryCode   string `json:"countryCode"`
	CountryName   string `json:"countryName"`
	CurrencyCode  string `json:"currencyCode"`
	GeonameID     int64  `json:"geonameId"`
	FipsCode      string `json:"fipsCode"`
	IsoAlpha3     string `json:"isoAlpha3"`
	IsoNumeric    string `json:"isoNumeric"`
	Languages     string `json:"languages"`
	Population    string `json:"population"`

	North float64 `json:"north"`
	South float64 `json:"south"`
	East  float64 `json:"east"`
	West  float64 `json:"west"`
}

// Country data model
// @Success 200 {object} countries.Country <-- This is a user defined struct.
type Country struct {
	AreaInSqKm    string   `json:"area_in_sq_km"`
	CallingCode   string   `json:"calling_code"`
	Capital       string   `json:"capital"`
	Continent     string   `json:"continent"`
	ContinentName string   `json:"continent_name"`
	CountryCode   string   `json:"country_code"`
	CountryName   string   `json:"country_name"`
	CurrencyCode  string   `json:"currency_code"`
	FipsCode      string   `json:"fips_code"`
	GeonameID     int64    `json:"geoname_id"`
	IsoAlpha3     string   `json:"iso_alpha3"`
	IsoNumeric    string   `json:"iso_numeric"`
	Locales       []string `json:"locales"`
	Population    string   `json:"population"`

	North float64 `json:"north"`
	South float64 `json:"south"`
	East  float64 `json:"east"`
	West  float64 `json:"west"`
}

// GeonamesGeoRes data model
type GeonamesGeoRes struct {
	TotalResultsCount int           `json:"totalResultsCount"`
	Geonames          []GeonamesGeo `json:"geonames"`
}

// GeonamesGeo data model
type GeonamesGeo struct {
	Name        string   `json:"name"`
	AdminCodes1 ISO31662 `json:"adminCodes1"`
	Lat         string   `json:"lat"`
	Lng         string   `json:"lng"`
	GeonameID   int64    `json:"geonameId"`
	CountryCode string   `json:"countryCode"`
	CountryName string   `json:"countryName"`
}

// ISO31662 data model
type ISO31662 struct {
	ISO31662 string `json:"ISO3166_2"`
}

// Geo data model
// @Success 200 {object} countries.Geo <-- This is a user defined struct.
type Geo struct {
	Name        string  `json:"name"`
	NameShort   string  `json:"name_short"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	GeonameID   int64   `json:"geoname_id"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
}

// NewCountriesFromGeonamesCountry creates a slice of countries from a slice of GeoNames countries sorted in country name order
func NewCountriesFromGeonamesCountry(ctx context.Context, geonamesCountries []GeonamesCountry) []Country {
	var countries []Country

	for _, geonamesCountry := range geonamesCountries {
		if _, ok := BlockListCountries[geonamesCountry.CountryCode]; ok {
			continue
		}

		callingCode := libphonenumber.GetCountryCodeForRegion(geonamesCountry.CountryCode)
		if callingCode == 0 {
			lib_log.Info(ctx, "Filtering out country due to error determining calling code for country", lib_log.FmtAny("geonamesCountry", geonamesCountry))
			continue
		}

		locales := strings.Split(geonamesCountry.Languages, ",")
		for i, v := range locales {
			locales[i] = strings.TrimSpace(v)
		}

		country := Country{
			AreaInSqKm:    geonamesCountry.AreaInSqKm,
			CallingCode:   fmt.Sprintf("%d", callingCode),
			Capital:       geonamesCountry.Capital,
			Continent:     geonamesCountry.Continent,
			ContinentName: geonamesCountry.ContinentName,
			CountryCode:   geonamesCountry.CountryCode,
			CountryName:   geonamesCountry.CountryName,
			CurrencyCode:  geonamesCountry.CurrencyCode,
			FipsCode:      geonamesCountry.FipsCode,
			GeonameID:     geonamesCountry.GeonameID,
			IsoAlpha3:     geonamesCountry.IsoAlpha3,
			IsoNumeric:    geonamesCountry.IsoNumeric,
			Locales:       locales,
			Population:    geonamesCountry.Population,

			North: geonamesCountry.North,
			South: geonamesCountry.South,
			East:  geonamesCountry.East,
			West:  geonamesCountry.West,
		}

		countries = append(countries, country)
	}

	sort.Slice(countries, func(i, j int) bool {
		return countries[i].CountryName < countries[j].CountryName
	})

	return countries
}

// PhoneNumberCountryCode returns the country calling code for the passed country code sorted in country name then name order
func PhoneNumberCountryCode(countryCode string) (string, error) {
	n, err := libphonenumber.Parse("6502530000", countryCode)
	if err != nil {
		return "", lib_errors.Wrap(err, "Failed parsing country code using libphonenumber")
	}
	return strconv.FormatInt(int64(n.GetCountryCode()), 10), nil
}

// NewGeosFromGeonamesGeos creates a slice of geos from a slice of GeoNames geos
func NewGeosFromGeonamesGeos(geonamesGeos []GeonamesGeo) ([]Geo, error) {
	var geos []Geo

	for _, geonamesGeo := range geonamesGeos {
		if _, ok := BlockListCountries[geonamesGeo.CountryCode]; ok {
			continue
		}

		lat, err := strconv.ParseFloat(geonamesGeo.Lat, 64)
		if err != nil {
			return nil, lib_errors.NewCustomf(http.StatusBadGateway, "Not recognized: geonamesGeo.Lat=%s", geonamesGeo.Lat)
		}

		lng, err := strconv.ParseFloat(geonamesGeo.Lng, 64)
		if err != nil {
			return nil, lib_errors.NewCustomf(http.StatusBadGateway, "Not recognized: geonamesGeo.Lng=%s", geonamesGeo.Lng)
		}

		geo := Geo{
			Name:        geonamesGeo.Name,
			NameShort:   geonamesGeo.AdminCodes1.ISO31662,
			Lat:         lat,
			Lng:         lng,
			GeonameID:   geonamesGeo.GeonameID,
			CountryCode: geonamesGeo.CountryCode,
			CountryName: geonamesGeo.CountryName,
		}

		sort.Slice(geos, func(i, j int) bool {
			if geos[i].CountryName == geos[j].CountryName {
				return geos[i].Name < geos[j].Name
			}
			return geos[i].CountryName < geos[j].CountryName
		})

		geos = append(geos, geo)
	}

	return geos, nil
}
