package svc

const (
	AccountsId             = "accounts-svc"
	AnalyticsId            = "analytics-svc"
	CmsId                  = "cms-svc"
	ContentfulSimilarityId = "contentful-similarity-svc"
	CoverId                = "cover-svc"
	CoverLcId              = "lc-api-cover"
	CustomerId             = "customer-svc"
	CustomerLcId           = "lc-api-customer"
	EmailId                = "email-svc"
	EngageClassifyId       = "engage-classify-svc"
	EngageId               = "engage-svc"
	IamId                  = "iam-svc"
	ImageId                = "image-svc"
	IssueId                = "issue-svc"
	ItemClassifyId         = "item-classify-svc"
	ItemId                 = "item-svc"
	ItemScrapeId           = "item-scrape-svc"
	LogisticsId            = "logistics-svc"
	OfacId                 = "ofac-svc"
	OrgId                  = "org-svc"
	PartnerId              = "partner-svc"
	PdfId                  = "pdf-svc"
	ProductClassifyId      = "product-classify-svc"
	ProductScrapeId        = "product-scrape-svc"
	ProductId              = "product-svc"
	ProgramId              = "program-svc"
	RateId                 = "rate-svc"
	RunId                  = "run-svc"
	SalesId                = "sales-svc"
	ScheduleId             = "schedule-svc"
	StorageId              = "storage-svc"
	TranslateId            = "translate-svc"
)

var (
	Services = []string{
		AccountsId,
		AnalyticsId,
		CmsId,
		ContentfulSimilarityId,
		CoverId,
		CoverLcId,
		CustomerId,
		CustomerLcId,
		EmailId,
		EngageClassifyId,
		EngageId,
		IamId,
		ImageId,
		IssueId,
		ItemClassifyId,
		ItemId,
		ItemScrapeId,
		LogisticsId,
		OfacId,
		OrgId,
		PartnerId,
		PdfId,
		ProductClassifyId,
		ProductScrapeId,
		ProductId,
		ProgramId,
		RateId,
		RunId,
		SalesId,
		ScheduleId,
		StorageId,
		TranslateId,
	}
)

func IsRecognizedId(svcId string) bool {
	switch svcId {
	default:
		return false
	case
		AccountsId,
		AnalyticsId,
		CmsId,
		ContentfulSimilarityId,
		CoverId,
		CoverLcId,
		CustomerId,
		CustomerLcId,
		EmailId,
		EngageClassifyId,
		EngageId,
		IamId,
		ImageId,
		IssueId,
		ItemClassifyId,
		ItemId,
		ItemScrapeId,
		LogisticsId,
		OfacId,
		OrgId,
		PartnerId,
		PdfId,
		ProductClassifyId,
		ProductScrapeId,
		ProductId,
		ProgramId,
		RateId,
		RunId,
		SalesId,
		ScheduleId,
		StorageId,
		TranslateId:
		return true
	}
}
