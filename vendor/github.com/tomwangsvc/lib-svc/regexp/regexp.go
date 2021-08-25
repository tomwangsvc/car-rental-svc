package regexp

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_svc "github.com/tomwangsvc/lib-svc/svc"
)

const (
	coverReferenceRegExp                     = `[0-9]{16}`
	coverReferenceContainedRegExp            = `(^|[^0-9])[0-9]{16}([^0-9]|$)`
	coverReferenceWithHyphensRegExp          = `[0-9]{4}\-[0-9]{4}\-[0-9]{4}\-[0-9]{4}`
	coverReferenceWithHyphensContainedRegExp = `(^|[^0-9])[0-9]{4}\-[0-9]{4}\-[0-9]{4}\-[0-9]{4}([^0-9]|$)`
	coverReferenceWithSpacesRegExp           = `[0-9]{4} [0-9]{4} [0-9]{4} [0-9]{4}`
	coverReferenceWithSpacesContainedRegExp  = `(^|[^0-9])[0-9]{4} [0-9]{4} [0-9]{4} [0-9]{4}([^0-9]|$)`
	claimReferenceRegExp                     = `CL10[0-9]{8}`
	claimReferenceContainedRegExp            = `CL10[0-9]{8}([^0-9]|$)`

	ConstantCaseRegExp          = `^[0-9A-Z_]+$` // TODO: test
	EmailRegExp                 = `^[^\s@]+@[^\s@]+\.[^\s@]+$`
	EmailWithOptionalNameRegExp = `^[^\s@]+@[^\s@]+\.[^\s@]+$|^.+<[^\s@]+@[^\s@]+\.[^\s@]+>$`
	MoneyRegExp                 = `^[0-9]+(\.[0-9]{1,2})?$`                                             // TODO: test
	WebsiteRegExp               = `^\S+\.\S+$`                                                          // TODO: test
	PhoneNumberRegExp           = `^[^a-zA-Z]+$`                                                        // TODO: test
	HourRegExp                  = `^([01][0-9]|2[0-3])(:([0-5][0-9]))?$`                                // TODO: test
	TimeZoneUtcOffsetExp        = `^-12:00|\+14:00|((-1[0-1]|-0[0-9]|\+0[0-9]|\+1[0-3]):([0-5][0-9]))$` // TODO: test
)

var (
	coverReferenceRegExpCompile                     *regexp.Regexp
	coverReferenceContainedRegExpCompile            *regexp.Regexp
	coverReferenceWithHyphensRegExpCompile          *regexp.Regexp
	coverReferenceWithHyphensContainedRegExpCompile *regexp.Regexp
	coverReferenceWithSpacesRegExpCompile           *regexp.Regexp
	coverReferenceWithSpacesContainedRegExpCompile  *regexp.Regexp
	claimReferenceRegExpCompile                     *regexp.Regexp
	claimReferenceContainedRegExpCompile            *regexp.Regexp
	timeZoneUtcOffsetExpCompile                     *regexp.Regexp

	internalUrlAppDomainRegExp     = `(goboxer|tomwang)\.com`
	internalUrlRegExpLocalhost     = `localhost:[0-9]+`
	internalUrlRegExpScheme        = `(https|http)://`
	internalUrlServiceDomainRegExp = `(goboxer\.cc|tomwang\.cc|tomwang\.net)`

	// We allow 'localhost' even in production environments since it is useful for testing
	// -> We only restrict public websites since browser restrictions apply to them, private websites (or even curl) could call us irrespective of any config we could apply
	InternalUrlAppDevRegExp = fmt.Sprintf(`^%s([^/]*\.dev\.%s|%s)(/.*|)$`, internalUrlRegExpScheme, internalUrlAppDomainRegExp, internalUrlRegExpLocalhost) // TODO: test
	InternalUrlAppPrdRegExp = fmt.Sprintf(`^%s([^/]*%s|%s)(/.*|)$`, internalUrlRegExpScheme, internalUrlAppDomainRegExp, internalUrlRegExpLocalhost)        // TODO: test
	InternalUrlAppStgRegExp = fmt.Sprintf(`^%s([^/]*\.stg\.%s|%s)(/.*|)$`, internalUrlRegExpScheme, internalUrlAppDomainRegExp, internalUrlRegExpLocalhost) // TODO: test
	InternalUrlAppUatRegExp = fmt.Sprintf(`^%s([^/]*\.uat\.%s|%s)(/.*|)$`, internalUrlRegExpScheme, internalUrlAppDomainRegExp, internalUrlRegExpLocalhost) // TODO: test

	InternalUrlTestAppRegExp = fmt.Sprintf(`^%s(test(\.dev|\.stg|\.uat|)\.%s|%s)(/.*|)$`, internalUrlRegExpScheme, internalUrlAppDomainRegExp, internalUrlRegExpLocalhost) // TODO: test

	InternalUrlBoxerAppDevRegExp = `^https://b49d5f4a1283d2ef.dev.(tomwang.com)(/.*|)$` // TODO: Change for go live, test
	InternalUrlBoxerAppPrdRegExp = `^https://b49d5f4a1283d2ef.prd.(tomwang.com)(/.*|)$` // TODO: Change for go live, test
	InternalUrlBoxerAppStgRegExp = `^https://b49d5f4a1283d2ef.stg.(tomwang.com)(/.*|)$` // TODO: Change for go live, test
	InternalUrlBoxerAppUatRegExp = `^https://b49d5f4a1283d2ef.uat.(tomwang.com)(/.*|)$` // TODO: Change for go live, test

	InternalUrlServiceRegExp string // TODO: test
)

func init() {
	var services []string
	for _, service := range lib_svc.Services {
		services = append(services, service)
		services = append(services, strings.Replace(strings.Replace(service, "lc-api-", "", -1), "-svc", "", -1))
		services = append(services, fmt.Sprintf("%s%s", strings.Replace((strings.Replace(service, "lc-api-", "", -1)), "-svc", "", -1), "lc-svc"))
	}
	internalServicesRegExp := fmt.Sprintf(`(%s)`, strings.Join(services, "|"))
	InternalUrlServiceRegExp = fmt.Sprintf(`^%s(%s\.(dev|prd|stg|uat)\.%s|%s)(/.*|)$`, internalUrlRegExpScheme, internalServicesRegExp, internalUrlServiceDomainRegExp, internalUrlRegExpLocalhost)

	coverReferenceRegExpCompile = regexp.MustCompile(coverReferenceRegExp)
	coverReferenceContainedRegExpCompile = regexp.MustCompile(coverReferenceContainedRegExp)
	coverReferenceWithHyphensRegExpCompile = regexp.MustCompile(coverReferenceWithHyphensRegExp)
	coverReferenceWithHyphensContainedRegExpCompile = regexp.MustCompile(coverReferenceWithHyphensContainedRegExp)
	coverReferenceWithSpacesRegExpCompile = regexp.MustCompile(coverReferenceWithSpacesRegExp)
	coverReferenceWithSpacesContainedRegExpCompile = regexp.MustCompile(coverReferenceWithSpacesContainedRegExp)
	claimReferenceRegExpCompile = regexp.MustCompile(claimReferenceRegExp)
	claimReferenceContainedRegExpCompile = regexp.MustCompile(claimReferenceContainedRegExp)
	timeZoneUtcOffsetExpCompile = regexp.MustCompile(TimeZoneUtcOffsetExp)
}

func FindCoverReference(ctx context.Context, s string) (coverReference string, ok bool) {
	lib_log.Info(ctx, "Finding", lib_log.FmtString("s", s))

	ok = coverReferenceContainedRegExpCompile.MatchString(s)
	if ok {
		coverReference = coverReferenceRegExpCompile.FindString(s)
	} else {
		ok = coverReferenceWithHyphensContainedRegExpCompile.MatchString(s)
		if ok {
			coverReference = strings.Replace(coverReferenceWithHyphensRegExpCompile.FindString(s), "-", "", -1)
		} else {
			ok = coverReferenceWithSpacesContainedRegExpCompile.MatchString(s)
			if ok {
				coverReference = strings.Replace(coverReferenceWithSpacesRegExpCompile.FindString(s), " ", "", -1)
			}
		}
	}

	lib_log.Info(ctx, "Found", lib_log.FmtString("coverReference", coverReference), lib_log.FmtBool("ok", ok))
	return
}

func FindClaimReference(ctx context.Context, s string) (claimReference string, ok bool) {
	lib_log.Info(ctx, "Finding", lib_log.FmtString("s", s))

	ok = claimReferenceContainedRegExpCompile.MatchString(s)
	if ok {
		claimReference = claimReferenceRegExpCompile.FindString(s)
	}

	lib_log.Info(ctx, "Found", lib_log.FmtString("claimReference", claimReference), lib_log.FmtBool("ok", ok))
	return
}

func ExtractHoursAndMinutesTimeZoneUtcOffset(timzoneUtcOffset string) (hours, minutes int64, err error) {
	if !timeZoneUtcOffsetExpCompile.MatchString(timzoneUtcOffset) {
		err = lib_errors.New("Cannot extract hours and minutes from given string")
		return
	}

	hours, err = strconv.ParseInt(timzoneUtcOffset[1:3], 10, 64)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed parsing string to int")
		return
	}
	minutes, err = strconv.ParseInt(timzoneUtcOffset[4:6], 10, 64)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed parsing string to int")
		return
	}

	if string(timzoneUtcOffset[0]) == "-" {
		hours *= -1
		minutes *= -1
	}
	return
}
