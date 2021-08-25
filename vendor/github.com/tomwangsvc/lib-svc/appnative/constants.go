package appnative

import (
	lib_env "github.com/tomwangsvc/lib-svc/env"
)

func isRecognizedCustomerAppAndroidPackageName(id, env string) bool {
	switch env {
	default:
		return false
	case lib_env.Dev:
		return id == "com.goboxer.boxer.dev"
	case lib_env.Prd:
		return id == "com.goboxer.boxer"
	case lib_env.Stg:
		return id == "com.goboxer.boxer.stg"
	case lib_env.Uat:
		return id == "com.goboxer.boxer.uat"
	}
}

func isRecognizedCustomerAppIosBundleId(id, env string) bool {
	switch env {
	default:
		return false
	case lib_env.Dev:
		return id == "com.goboxer.BoxerDev" || id == "com.goboxer.BoxerDevDebug"
	case lib_env.Prd:
		return id == "com.goboxer.BoxerPrd" || id == "com.goboxer.Boxer"
	case lib_env.Stg:
		return id == "com.goboxer.BoxerStg"
	case lib_env.Uat:
		return id == "com.goboxer.BoxerUat"
	}
}

func IsRecognizedId(id, env string) bool {
	return isRecognizedCustomerAppAndroidPackageName(id, env) || isRecognizedCustomerAppIosBundleId(id, env)
}
