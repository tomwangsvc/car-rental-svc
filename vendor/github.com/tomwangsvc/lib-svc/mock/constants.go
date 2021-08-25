package mock

import (
	"net/http"
	"time"

	lib_countries "github.com/tomwangsvc/lib-svc/countries"
	lib_http "github.com/tomwangsvc/lib-svc/http"
)

var (
	ExpectedResultBytes   = []byte("EXPECTED")
	ExpectedResultBool    = true
	ExpectedResultFloat64 = 1.0
	ExpectedResultInt     = 1
	ExpectedResultInt64   = int64(ExpectedResultInt)
	ExpectedResultString  = "EXPECTED"
	ExpectedResultStrings = []string{ExpectedResultString}
	ExpectedResultTime    = time.Time{}.Add(1)

	ExpectedHeaderForXLcPaginationTotalOfZero = http.Header{lib_http.HeaderKeyXLcPaginationTotal: []string{"0"}}
	ExpectedHeaderForXLcPaginationTotalOfOne  = http.Header{lib_http.HeaderKeyXLcPaginationTotal: []string{"1"}}

	ExpectedLibCountry                = lib_countries.Country{}
	ExpectedLibCountriesByCountryCode = map[string]lib_countries.Country{
		ExpectedResultString: {},
	}
	ExpectedLibCountriesByCurrencyCode = map[string][]lib_countries.Country{
		ExpectedResultString: {},
	}
	ExpectedResultLibCountryMetadata = lib_countries.Metadata{
		CountriesByCountryCode:  ExpectedLibCountriesByCountryCode,
		CountriesByCurrencyCode: ExpectedLibCountriesByCurrencyCode,
	}

	TokenForLibTokenGcpVerifyWithBadGateway                = "TOKEN_FOR_LIB_TOKEN_GCP_VERIFY_WITH_BAD_GATEWAY"
	TokenForLibTokenGcpVerifyWithGatewayTimeout            = "TOKEN_FOR_LIB_TOKEN_GCP_VERIFY_WITH_GATEWAY_TIMEOUT"
	TokenForLibTokenGcpVerifyWithServiceUnavailable        = "TOKEN_FOR_LIB_TOKEN_GCP_VERIFY_WITH_SERVICE_UNAVAILABLE"
	TokenForLibTokenIamVerifyAndExtractClaimsWithTestFalse = "TOKEN_FOR_LIB_TOKEN_IAM_VERIFY_AND_EXTRACT_CLAIMS_WITH_TEST_FALSE"
)
