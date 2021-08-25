package brand

import (
	"fmt"

	lib_app "github.com/tomwangsvc/lib-svc/app"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
)

const (
	BoxerEmailDomain      = "goboxer.com"
	BoxerId               = "BOXER"
	BoxerIdentityProvider = "cicp-customer-app"
	BoxerName             = "Boxer"

	TomWangEmailDomain      = "tomwang.com"
	TomWangId               = "TOMWANG"
	TomWangIdentityProvider = "cicp-customerlc-app"
	TomWangName             = "TomWang" // We deliberately use a generic name so it works for all TomWang companies registered in all territories
	TomWangUuid             = "578a5603-dd1a-4d03-b570-d6cf04c8aee6"
)

func IsBoxerIdentityProvider(identityProvider string) bool {
	return identityProvider == BoxerIdentityProvider
}

func IsTomWangIdentityProvider(identityProvider string) bool {
	return identityProvider == TomWangIdentityProvider
}

func EmailDomain(id string) (string, error) {
	switch id {
	default:
		return "", lib_errors.Errorf("Brand ID %q not recognized", id)
	case BoxerId:
		return BoxerEmailDomain, nil
	case TomWangId:
		return TomWangEmailDomain, nil
	}
}

func Name(id string) (string, error) {
	switch id {
	default:
		return "", lib_errors.Errorf("Brand ID %q not recognized", id)
	case BoxerId:
		return BoxerName, nil
	case TomWangId:
		return TomWangName, nil
	}
}

func Recognized(id string) bool {
	switch id {
	default:
		return false
	case
		BoxerId,
		TomWangId:
		return true
	}
}

func LinkBoxerHelp(envId string) string {
	return fmt.Sprintf("%s/help", lib_app.CustomerWebAppBaseUrl(envId))
}

func LinkBoxerCustomerAppWebModeMarketing(envId string) string {
	return fmt.Sprintf("%s/?mode=marketing", lib_app.CustomerWebAppBaseUrl(envId))
}

func LinkBoxerPrivacyPolicy(envId string) string {
	switch envId {
	default:
		// return "https://dev.goboxer.com/privacy" TODO: Use this domain when available
		return "https://l180922426832059.web.app/privacy"
	case lib_env.Dev:
		// return "https://dev.goboxer.com/privacy" TODO: Use this domain when available
		return "https://l180922426832059.web.app/privacy"
	case lib_env.Prd:
		// return "https://goboxer.com/privacy" TODO: Use this domain when available
		return "https://l180922426832059.web.app/privacy"
	case lib_env.Stg:
		// return "https://stg.goboxer.com/privacy" TODO: Use this domain when available
		return "https://l180922426832059.web.app/privacy"
	case lib_env.Uat:
		// return "https://uat.goboxer.com/privacy" TODO: Use this domain when available
		return "https://l180922426832059.web.app/privacy"
	}
}

func LinkBoxerContact(envId string) string {
	switch envId {
	default:
		// return "https://dev.goboxer.com/contact" TODO: Use this domain when available
		return "https://l180922426832059.web.app/contact"
	case lib_env.Dev:
		// return "https://dev.goboxer.com/contact" TODO: Use this domain when available
		return "https://l180922426832059.web.app/contact"
	case lib_env.Prd:
		// return "https://goboxer.com/contact" TODO: Use this domain when available
		return "https://l180922426832059.web.app/contact"
	case lib_env.Stg:
		// return "https://stg.goboxer.com/contact" TODO: Use this domain when available
		return "https://l180922426832059.web.app/contact"
	case lib_env.Uat:
		// return "https://uat.goboxer.com/contact" TODO: Use this domain when available
		return "https://l180922426832059.web.app/contact"
	}
}

// TODO: considering returning error
func LinkBoxerGetInTouch(envId string) string {
	switch envId {
	default:
		// return "https://dev.goboxer.com/get-in-touch" TODO: Use this domain when available
		return "https://l180922426832059.web.app/get-in-touch"
	case lib_env.Dev:
		// return "https://dev.goboxer.com/get-in-touch" TODO: Use this domain when available
		return "https://l180922426832059.web.app/get-in-touch"
	case lib_env.Prd:
		// return "https://goboxer.com/get-in-touch" TODO: Use this domain when available
		return "https://l180922426832059.web.app/get-in-touch"
	case lib_env.Stg:
		// return "https://stg.goboxer.com/get-in-touch" TODO: Use this domain when available
		return "https://l180922426832059.web.app/get-in-touch"
	case lib_env.Uat:
		// return "https://uat.goboxer.com/get-in-touch" TODO: Use this domain when available
		return "https://l180922426832059.web.app/get-in-touch"
	}
}

func LinkBoxerWhatBoxerIs(envId string) string {
	switch envId {
	default:
		// return "https://dev.goboxer.com/what-boxer-is" TODO: Use this domain when available
		return "https://l180922426832059.web.app/what-boxer-is"
	case lib_env.Dev:
		// return "https://dev.goboxer.com/what-boxer-is" TODO: Use this domain when available
		return "https://l180922426832059.web.app/what-boxer-is"
	case lib_env.Prd:
		// return "https://goboxer.com/what-boxer-is" TODO: Use this domain when available
		return "https://l180922426832059.web.app/what-boxer-is"
	case lib_env.Stg:
		// return "https://stg.goboxer.com/what-boxer-is" TODO: Use this domain when available
		return "https://l180922426832059.web.app/what-boxer-is"
	case lib_env.Uat:
		// return "https://uat.goboxer.com/what-boxer-is" TODO: Use this domain when available
		return "https://l180922426832059.web.app/what-boxer-is"
	}
}

func LinkBoxerHelpSupportYourCoverCategory(envId string) string {
	return fmt.Sprintf("%s/help/category:your_covers", lib_app.CustomerWebAppBaseUrl(envId))
}

func LinkBoxerSalesAppWeb(envId string) string {
	switch envId {
	default:
		// return "https://dev.goboxer.com" TODO: Use this domain when available
		return "https://l180922426832059.web.app"
	case lib_env.Dev:
		// return "https://dev.goboxer.com" TODO: Use this domain when available
		return "https://l180922426832059.web.app"
	case lib_env.Prd:
		// return "https://goboxer.com" TODO: Use this domain when available
		return "https://l180922426832059.web.app"
	case lib_env.Stg:
		// return "https://stg.goboxer.com" TODO: Use this domain when available
		return "https://l180922426832059.web.app"
	case lib_env.Uat:
		// return "https://uat.goboxer.com" TODO: Use this domain when available
		return "https://l180922426832059.web.app"
	}
}

func LinkBoxerMobileApp(envId string) string {
	switch envId {
	default:
		return "https://dev.goboxer.com/app"
	case lib_env.Dev:
		return "https://dev.goboxer.com/app"
	case lib_env.Prd:
		return "https://goboxer.com/app"
	case lib_env.Stg:
		return "https://stg.goboxer.com/app"
	case lib_env.Uat:
		return "https://uat.goboxer.com/app"
	}
}

func LinkBoxerCoverDetailsWithCoverId(coverId, envId string) string {
	return fmt.Sprintf("%s/coverdetails/%s", lib_app.CustomerWebAppBaseUrl(envId), coverId)
}

func LinkBoxerForLoginWithToken(env, token string) string {
	return fmt.Sprintf("%s/?token=%s", lib_app.CustomerWebAppBaseUrl(env), token)
}
