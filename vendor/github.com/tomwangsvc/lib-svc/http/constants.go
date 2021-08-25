package http

const (
	// ---------- Non-Standard Header keys & values ----------
	HeaderKeyXLcCorrelationId = "X-Lc-Correlation-Id"

	HeaderKeyXLcCloudTaskCreatedDate = "X-Lc-Cloud-Task-Created-Date"
	HeaderKeyXLcCloudTaskId          = "X-Lc-Cloud-Task-Id"

	HeaderKeyXLcCaller = "X-Lc-Caller"

	HeaderKeyXLcETagsForObjects = "X-Lc-ETags-For-Objects"

	// Added to requests by the GCLB
	HeaderKeyXLcLocationCity              = "X-Lc-Location-City"
	HeaderKeyXLcLocationCityLatLng        = "X-Lc-Location-City-Lat-Lng"
	HeaderKeyXLcLocationCountry           = "X-Lc-Location-Country" // TODO - Remove this ASAP @2021/02, it has been replaced by HeaderKeyXLcLocationRegion & HeaderKeyXLcLocationRegionSubdivision due to now we have https://cloud.google.com/load-balancing/docs/custom-headers
	HeaderKeyXLcLocationIp                = "X-Lc-Location-Ip"
	HeaderKeyXLcLocationRegion            = "X-Lc-Location-Region"
	HeaderKeyXLcLocationRegionSubdivision = "X-Lc-Location-Region-Subdivision"

	HeaderKeyXLcPaginationCursor        = "X-Lc-Pagination-Cursor"
	HeaderKeyXLcPaginationLimit         = "X-Lc-Pagination-Limit"
	HeaderKeyXLcPaginationOffset        = "X-Lc-Pagination-Offset"
	HeaderKeyXLcPaginationOrder         = "X-Lc-Pagination-Order"
	HeaderKeyXLcPaginationReadTimestamp = "X-Lc-Pagination-Read-Timestamp"
	HeaderKeyXLcPaginationTotal         = "X-Lc-Pagination-Total"

	HeaderKeyXLcSvcIntegrationTest   = "X-Lc-Svc-Integration-Test"
	HeaderValueXLcSvcIntegrationTest = "0aaee093-8b4a-48cc-8d5a-c08559ccd8fb"

	HeaderKeyXLcSvcIntegrationTestPubsubAutoAckDisable = "X-Lc-Svc-Integration-Test-Pubsub-Auto-Ack-Disable"

	HeaderKeyXLcPubsubMessageId          = "X-Lc-Pubsub-Message-Id"
	HeaderKeyXLcPubsubMessagePublishTime = "X-Lc-Pubsub-Message-Publish-Time"

	HeaderKeyXLcSvcTest   = "X-Lc-Svc-Test"
	HeaderValueXLcSvcTest = "bdfcdc57-b05a-49c6-8b1f-cfd01272411e"

	// Used to pass metadata to clients associated with an api running in test mode
	// -> E.g. Short circuit email delivery of validation tokens by instead returning them in a test header
	HeaderKeyXLcSvcTestMetadata = "X-Lc-Svc-Test-Metadata"

	// Clients do not use the scheme when sending this origin
	xLcLocalhostCallerHeader = "localhost"
)

func XLcLocationHeader(headerKey string) bool {
	return headerKey == HeaderKeyXLcLocationCity ||
		headerKey == HeaderKeyXLcLocationCityLatLng ||
		headerKey == HeaderKeyXLcLocationCountry ||
		headerKey == HeaderKeyXLcLocationIp ||
		headerKey == HeaderKeyXLcLocationRegion ||
		headerKey == HeaderKeyXLcLocationRegionSubdivision
}
