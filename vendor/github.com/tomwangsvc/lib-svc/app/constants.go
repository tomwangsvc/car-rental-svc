package app

import (
	"fmt"
	"regexp"
	"strings"

	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_regexp "github.com/tomwangsvc/lib-svc/regexp"
)

const (
	CicpId              = "cicp-app"
	CicpDevGcpProjectId = "tomwang-100"
	CicpPrdGcpProjectId = "tomwang-102"
	CicpStgGcpProjectId = "l291651910754051"
	CicpUatGcpProjectId = "tomwang-101"

	CicpCustomerId              = "cicp-customer-app"
	CicpCustomerDevGcpProjectId = "l194333055636611"
	CicpCustomerPrdGcpProjectId = "l173842553566022"
	CicpCustomerStgGcpProjectId = "l948429608226318"
	CicpCustomerUatGcpProjectId = "l184082804621515"

	CicpCustomerLcId              = "cicp-customerlc-app"
	CicpCustomerLcDevGcpProjectId = "l216261647120632"
	CicpCustomerLcPrdGcpProjectId = "l115731743383503"
	CicpCustomerLcStgGcpProjectId = "l213543134416847"
	CicpCustomerLcUatGcpProjectId = "l110902466012751"

	CicpPartnerId              = "cicp-partner-app"
	CicpPartnerDevGcpProjectId = "l103422215933782"
	CicpPartnerPrdGcpProjectId = "l239729629715014"
	CicpPartnerStgGcpProjectId = "l172951737111189"
	CicpPartnerUatGcpProjectId = "l257891232876532"

	CmsId              = "cms-app"
	CmsDevGcpProjectId = "l119610946166937"
	CmsPrdGcpProjectId = "l138591287029267"
	CmsStgGcpProjectId = "l167490771305160"
	CmsUatGcpProjectId = "l239111190822980"

	CrmId              = "crm-app"
	CrmDevGcpProjectId = "tomwang-35"
	CrmPrdGcpProjectId = "tomwang-64"
	CrmStgGcpProjectId = "l455623310216973"
	CrmUatGcpProjectId = "tomwang-63"

	CustomerId              = "customer-app"
	CustomerDevGcpProjectId = "l304832640514857"
	CustomerPrdGcpProjectId = "l433329581903530"
	CustomerStgGcpProjectId = "l237911610530099"
	CustomerUatGcpProjectId = "l184451596918406"

	CustomerLcId              = "customerlc-app"
	CustomerLcDevGcpProjectId = "tomwang-39"
	CustomerLcPrdGcpProjectId = "tomwang-5"
	CustomerLcStgGcpProjectId = "l328426628405415"
	CustomerLcUatGcpProjectId = "tomwang-6"

	EngageId              = "engage-app"
	EngageDevGcpProjectId = "l266041656721199"
	EngagePrdGcpProjectId = "l601525719171861"
	EngageStgGcpProjectId = "l147102962238197"
	EngageUatGcpProjectId = "l178852745511402"

	PartnerId              = "partner-app"
	PartnerDevGcpProjectId = "l284374898209082"
	PartnerPrdGcpProjectId = "l179011308731595"
	PartnerStgGcpProjectId = "l184681824524846"
	PartnerUatGcpProjectId = "l133586341303387"
)

var (
	compiledInternalUrlTestAppRegExp = regexp.MustCompile(lib_regexp.InternalUrlTestAppRegExp)
)

func AllIds() []string {
	ids := []string{
		CicpId,
		CicpCustomerId,
		CicpCustomerLcId,
		CicpPartnerId,
		CmsId,
		CrmId,
		CustomerId,
		CustomerLcId,
		EngageId,
		PartnerId,
	}
	return ids
}

func AllIdsWithEnvs(_ string) []string {
	return []string{}
}

func IsPublicId(appId string) bool {
	switch appId {
	default:
		return false
	case CustomerLcId, CustomerId, PartnerId:
		return true
	}
}

func IsTestAppUrl(url string) bool {
	return compiledInternalUrlTestAppRegExp.Match([]byte(url))
}

func GcpProjectId(_, _ string) string {
	return ""
}

//revive:disable:cyclomatic
func Id(gcpProjectId string) string {
	switch gcpProjectId {
	default:
		return ""

	case CicpDevGcpProjectId,
		CicpPrdGcpProjectId,
		CicpStgGcpProjectId,
		CicpUatGcpProjectId:
		return CicpId

	case CicpCustomerDevGcpProjectId,
		CicpCustomerPrdGcpProjectId,
		CicpCustomerStgGcpProjectId,
		CicpCustomerUatGcpProjectId:
		return CicpCustomerId

	case CicpCustomerLcDevGcpProjectId,
		CicpCustomerLcPrdGcpProjectId,
		CicpCustomerLcStgGcpProjectId,
		CicpCustomerLcUatGcpProjectId:
		return CicpCustomerLcId

	case CicpPartnerDevGcpProjectId,
		CicpPartnerPrdGcpProjectId,
		CicpPartnerStgGcpProjectId,
		CicpPartnerUatGcpProjectId:
		return CicpPartnerId

	case CmsDevGcpProjectId,
		CmsPrdGcpProjectId,
		CmsStgGcpProjectId,
		CmsUatGcpProjectId:
		return CmsId

	case CrmDevGcpProjectId,
		CrmPrdGcpProjectId,
		CrmStgGcpProjectId,
		CrmUatGcpProjectId:
		return CrmId

	case CustomerDevGcpProjectId,
		CustomerPrdGcpProjectId,
		CustomerStgGcpProjectId,
		CustomerUatGcpProjectId:
		return CustomerId

	case CustomerLcDevGcpProjectId,
		CustomerLcPrdGcpProjectId,
		CustomerLcStgGcpProjectId,
		CustomerLcUatGcpProjectId:
		return CustomerLcId

	case EngageDevGcpProjectId,
		EngagePrdGcpProjectId,
		EngageStgGcpProjectId,
		EngageUatGcpProjectId:
		return EngageId

	case PartnerDevGcpProjectId,
		PartnerPrdGcpProjectId,
		PartnerStgGcpProjectId,
		PartnerUatGcpProjectId:
		return PartnerId
	}
	//revive:enable:cyclomatic
}

//revive:disable:cyclomatic
func IdWithEnv(_, _ string) string {
	return ""
	//revive:enable:cyclomatic
}

//revive:disable:cyclomatic
func IdWithEnvForGcpProjectId(gcpProjectId string) string {
	switch gcpProjectId {
	default:
		return ""

	case CicpDevGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpId, lib_env.Dev)
	case CicpPrdGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpId, lib_env.Prd)
	case CicpStgGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpId, lib_env.Stg)
	case CicpUatGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpId, lib_env.Uat)

	case CicpCustomerDevGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpCustomerId, lib_env.Dev)
	case CicpCustomerPrdGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpCustomerId, lib_env.Prd)
	case CicpCustomerStgGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpCustomerId, lib_env.Stg)
	case CicpCustomerUatGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpCustomerId, lib_env.Uat)

	case CicpCustomerLcDevGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpCustomerLcId, lib_env.Dev)
	case CicpCustomerLcPrdGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpCustomerLcId, lib_env.Prd)
	case CicpCustomerLcStgGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpCustomerLcId, lib_env.Stg)
	case CicpCustomerLcUatGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpCustomerLcId, lib_env.Uat)

	case CicpPartnerDevGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpPartnerId, lib_env.Dev)
	case CicpPartnerPrdGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpPartnerId, lib_env.Prd)
	case CicpPartnerStgGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpPartnerId, lib_env.Stg)
	case CicpPartnerUatGcpProjectId:
		return fmt.Sprintf("%s-%s", CicpPartnerId, lib_env.Uat)

	case CmsDevGcpProjectId:
		return fmt.Sprintf("%s-%s", CmsId, lib_env.Dev)
	case CmsPrdGcpProjectId:
		return fmt.Sprintf("%s-%s", CmsId, lib_env.Prd)
	case CmsStgGcpProjectId:
		return fmt.Sprintf("%s-%s", CmsId, lib_env.Stg)
	case CmsUatGcpProjectId:
		return fmt.Sprintf("%s-%s", CmsId, lib_env.Uat)

	case CrmDevGcpProjectId:
		return fmt.Sprintf("%s-%s", CrmId, lib_env.Dev)
	case CrmPrdGcpProjectId:
		return fmt.Sprintf("%s-%s", CrmId, lib_env.Prd)
	case CrmStgGcpProjectId:
		return fmt.Sprintf("%s-%s", CrmId, lib_env.Stg)
	case CrmUatGcpProjectId:
		return fmt.Sprintf("%s-%s", CrmId, lib_env.Uat)

	case CustomerDevGcpProjectId:
		return fmt.Sprintf("%s-%s", CustomerId, lib_env.Dev)
	case CustomerPrdGcpProjectId:
		return fmt.Sprintf("%s-%s", CustomerId, lib_env.Prd)
	case CustomerStgGcpProjectId:
		return fmt.Sprintf("%s-%s", CustomerId, lib_env.Stg)
	case CustomerUatGcpProjectId:
		return fmt.Sprintf("%s-%s", CustomerId, lib_env.Uat)

	case CustomerLcDevGcpProjectId:
		return fmt.Sprintf("%s-%s", CustomerLcId, lib_env.Dev)
	case CustomerLcPrdGcpProjectId:
		return fmt.Sprintf("%s-%s", CustomerLcId, lib_env.Prd)
	case CustomerLcStgGcpProjectId:
		return fmt.Sprintf("%s-%s", CustomerLcId, lib_env.Stg)
	case CustomerLcUatGcpProjectId:
		return fmt.Sprintf("%s-%s", CustomerLcId, lib_env.Uat)

	case EngageDevGcpProjectId:
		return fmt.Sprintf("%s-%s", EngageId, lib_env.Dev)
	case EngagePrdGcpProjectId:
		return fmt.Sprintf("%s-%s", EngageId, lib_env.Prd)
	case EngageStgGcpProjectId:
		return fmt.Sprintf("%s-%s", EngageId, lib_env.Stg)
	case EngageUatGcpProjectId:
		return fmt.Sprintf("%s-%s", EngageId, lib_env.Uat)

	case PartnerDevGcpProjectId:
		return fmt.Sprintf("%s-%s", PartnerId, lib_env.Dev)
	case PartnerPrdGcpProjectId:
		return fmt.Sprintf("%s-%s", PartnerId, lib_env.Prd)
	case PartnerStgGcpProjectId:
		return fmt.Sprintf("%s-%s", PartnerId, lib_env.Stg)
	case PartnerUatGcpProjectId:
		return fmt.Sprintf("%s-%s", PartnerId, lib_env.Uat)
	}
	//revive:enable:cyclomatic
}

func CheckEnvOfIdWithEnv(appIdWithEnv, env string) error {
	switch env {
	default:
		return lib_errors.Errorf("Environment not recognized %q", env)
	case lib_env.Dev:
		if !strings.HasSuffix(appIdWithEnv, lib_env.Dev) {
			return lib_errors.Errorf("Invalid suffix for appIdWithEnv %q, expected %q", appIdWithEnv, lib_env.Dev)
		}
	case lib_env.Prd:
		if !strings.HasSuffix(appIdWithEnv, lib_env.Prd) {
			return lib_errors.Errorf("Invalid suffix for appIdWithEnv %q, expected %q", appIdWithEnv, lib_env.Prd)
		}
	case lib_env.Stg:
		if !strings.HasSuffix(appIdWithEnv, lib_env.Stg) {
			return lib_errors.Errorf("Invalid suffix for appIdWithEnv %q, expected %q", appIdWithEnv, lib_env.Stg)
		}
	case lib_env.Uat:
		if !strings.HasSuffix(appIdWithEnv, lib_env.Uat) {
			return lib_errors.Errorf("Invalid suffix for appIdWithEnv %q, expected %q", appIdWithEnv, lib_env.Uat)
		}
	}
	return nil
}

// TODO - Remove/Change this when the domain goes live
func CustomerWebAppBaseUrl(env string) string {
	switch env {
	case lib_env.Prd:
		return fmt.Sprintf("https://foo.%s.tomwang.com", lib_env.Prd)
	case lib_env.Stg:
		return fmt.Sprintf("https://foo.%s.tomwang.com", lib_env.Stg)
	case lib_env.Uat:
		return fmt.Sprintf("https://foo.%s.tomwang.com", lib_env.Uat)
	default:
		return fmt.Sprintf("https://foo.%s.tomwang.com", lib_env.Dev)
	}
}
