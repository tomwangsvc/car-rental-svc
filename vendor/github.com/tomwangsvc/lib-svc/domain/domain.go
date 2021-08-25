package domain

import (
	lib_strings "github.com/tomwangsvc/lib-svc/strings"
)

const (
	Accountant              = "ACCOUNTANT"
	AccountsPaymentBatch    = "ACCOUNTS_PAYMENT_BATCH"
	AccountsPaymentInvoice  = "ACCOUNTS_PAYMENT_INVOICE"
	AccountsSvc             = "ACCOUNTS_SVC"
	AnalyticsSvc            = "ANALYTICS_SVC"
	Boxer                   = "BOXER"
	Claim                   = "CLAIM"
	ClaimAction             = "CLAIM_ACTION"
	ClaimIssue              = "CLAIM_ISSUE"
	ClaimIssueCustomer      = "CLAIM_ISSUE_CUSTOMER"
	ClaimJob                = "CLAIM_JOB"
	ClaimSvc                = "CLAIM_SVC"
	ClaimTask               = "CLAIM_TASK"
	ClaimTaskCustomer       = "CLAIM_TASK_CUSTOMER"
	ClaimWorkflow           = "CLAIM_WORKFLOW"
	ContentfulSimilaritySvc = "CONTENTFUL_SIMULARITY_SVC"
	Cover                   = "COVER"
	CoverCustomer           = "COVER_CUSTOMER"
	CoverLcSvc              = "COVERLC_SVC"
	CoverSvc                = "COVER_SVC"
	Csr                     = "CSR"
	CsrNote                 = "CSR_NOTE"
	Customer                = "CUSTOMER"
	CustomerAccount         = "CUSTOMER_ACCOUNT"
	CustomerLcSvc           = "CUSTOMERLC_SVC"
	CustomerSvc             = "CUSTOMER_SVC"
	Document                = "DOCUMENT"
	Email                   = "EMAIL"
	EmailIn                 = "EMAIL_IN"
	EmailOut                = "EMAIL_OUT"
	EmailSvc                = "EMAIL_SVC"
	EngageMessage           = "ENGAGE_MESSAGE"
	EngageMessageChannel    = "ENGAGE_MESSAGE_CHANNEL"
	EngageSvc               = "ENGAGE_SVC"
	IamSvc                  = "IAM_SVC"
	TomWang                 = "TOMWANG"
	LogisticsShipment       = "LOGISTICS_SHIPMENT"
	LogisticsSvc            = "LOGISTICS_SVC"
	Manufacturer            = "MANUFACTURER"
	Marketing               = "MARKETING"
	Merchant                = "MERCHANT"
	MerchantSvc             = "MERCHANT_SVC"
	MlClaim                 = "ML_CLAIM"
	MlCover                 = "ML_COVER"
	MlCustomerAccount       = "ML_CUSTOMER_ACCOUNT"
	Org                     = "ORG"
	OrgCarrier              = "ORG_CARRIER"
	OrgLogistics            = "ORG_LOGISTICS"
	OrgManufacturer         = "ORG_MANUFACTURER"
	OrgMerchant             = "ORG_MERCHANT"
	OrgProvider             = "ORG_PROVIDER"
	OrgRepairer             = "ORG_REPAIRER"
	OrgSeller               = "ORG_SELLER"
	OrgSvc                  = "ORG_SVC"
	OrgWarehouse            = "ORG_WAREHOUSE"
	ProductBrand            = "PRODUCT_BRAND"
	ProductCategory         = "PRODUCT_CATEGORY"
	ProductClass            = "PRODUCT_CLASS"
	ProductClassifyBetaSvc  = "PRODUCT_CLASSIFY_BETA_SVC"
	ProductClassifySvc      = "PRODUCT_CLASSIFY_SVC"
	ProductModel            = "PRODUCT_MODEL"
	ProductSeries           = "PRODUCT_SERIES"
	ProductSubcategory      = "PRODUCT_SUBCATEGORY"
	ProductSvc              = "PRODUCT_SVC"
	Repairer                = "REPAIRER"
	RepairerSvc             = "REPAIRER_SVC"
	SalesSvc                = "SALES_SVC"
	ScheduleSvc             = "SCHEDULE_SVC"
	Seller                  = "SELLER"
	StorageSvc              = "STORAGE_SVC"
	System                  = "SYSTEM"
	TestAppWeb              = "TEST_APP_WEB"
	Unknown                 = "UNKNOWN"
)

var (
	FromOrTo = []string{
		Accountant,
		Boxer,
		Csr,
		Customer,
		CustomerAccount,
		TomWang,
		Org,
		Unknown,
	}

	OrgSpecializations = []string{
		OrgCarrier,
		OrgLogistics,
		OrgManufacturer,
		OrgMerchant,
		OrgProvider,
		OrgRepairer,
		OrgSeller,
		OrgWarehouse,
	}
	Orgs = append(OrgSpecializations, Org)

	RecognizedDomains = []string{
		Accountant,
		AccountsPaymentBatch,
		AccountsPaymentInvoice,
		AccountsSvc,
		AnalyticsSvc,
		Boxer,
		Claim,
		ClaimAction,
		ClaimIssue,
		ClaimIssueCustomer,
		ClaimJob,
		ClaimSvc,
		ClaimTask,
		ClaimTaskCustomer,
		ClaimWorkflow,
		ContentfulSimilaritySvc,
		Cover,
		CoverCustomer,
		CoverLcSvc,
		CoverSvc,
		Csr,
		CsrNote,
		Customer,
		CustomerAccount,
		CustomerLcSvc,
		CustomerSvc,
		Document,
		Email,
		EmailIn,
		EmailOut,
		EmailSvc,
		EngageMessage,
		EngageMessageChannel,
		EngageSvc,
		IamSvc,
		TomWang,
		LogisticsShipment,
		LogisticsSvc,
		Manufacturer,
		Marketing,
		Merchant,
		MerchantSvc,
		MlClaim,
		MlCover,
		MlCustomerAccount,
		Org,
		OrgCarrier,
		OrgLogistics,
		OrgManufacturer,
		OrgMerchant,
		OrgProvider,
		OrgRepairer,
		OrgSeller,
		OrgSvc,
		OrgWarehouse,
		ProductBrand,
		ProductCategory,
		ProductClass,
		ProductClassifyBetaSvc,
		ProductClassifySvc,
		ProductModel,
		ProductSeries,
		ProductSubcategory,
		ProductSvc,
		Repairer,
		RepairerSvc,
		SalesSvc,
		ScheduleSvc,
		Seller,
		StorageSvc,
		System,
		TestAppWeb,
		Unknown,
	}
)

func Recognized(domain string) bool {
	return lib_strings.Contains(RecognizedDomains, domain)
}

func RecognizedFromOrTo(domain string) bool {
	return lib_strings.Contains(FromOrTo, domain)
}

func RecognizedOrg(domain string) bool {
	return lib_strings.Contains(Orgs, domain)
}

func RecognizedOrgSpecialization(domain string) bool {
	return lib_strings.Contains(OrgSpecializations, domain)
}
